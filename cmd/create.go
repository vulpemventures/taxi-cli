package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/tiero/ocean/pkg/coinselect"
	"github.com/tiero/ocean/pkg/confidential"
	"github.com/tiero/ocean/pkg/explorer/blockstream"
	"github.com/tiero/ocean/pkg/keypair"
	"github.com/tiero/ocean/pkg/partial"

	"github.com/vulpemventures/go-elements/address"
	"github.com/vulpemventures/go-elements/network"
	"github.com/vulpemventures/taxi-protobuf/rpc/taxi"
)

func createAction(from string, to string, asset string, amount int, blinding string, regtest bool, explorerURL string, taxiURL string) error {
	if from == "" || to == "" || asset == "" || amount <= 0 {
		return fmt.Errorf("missing required flag")
	}

	// Explorer
	e := blockstream.NewExplorer(explorerURL)
	// Network
	currentNetwork := network.Liquid
	if regtest {
		currentNetwork = network.Regtest
	}

	fromScript, err := address.ToOutputScript(from, currentNetwork)
	if err != nil {
		return fmt.Errorf("from: %w", err)
	}
	toScript, err := address.ToOutputScript(to, currentNetwork)
	if err != nil {
		return fmt.Errorf("to: %w", err)
	}

	var bk []byte
	var isConfidential = false
	if len(blinding) > 0 {
		isConfidential = true
		bk, err = hex.DecodeString(blinding)
		if err != nil {
			return fmt.Errorf("Blinding Key: %w", err)
		}
	}

	utxos, err := e.GetUnspents(from)
	if err != nil {
		return fmt.Errorf("unspents: %w", err)
	}

	coins := &coinselect.Coins{Utxos: utxos}
	if isConfidential {
		coins.BlindingKey = bk
	}
	selectedUtxos, change, err := coins.CoinSelect(uint64(amount), asset)
	if err != nil {
		return fmt.Errorf("coin selection: %w", err)
	}

	psetWithoutFees := partial.NewPartial(&currentNetwork)
	blindingPrivKeysOfInputs := make([][]byte, 0)
	for _, utxo := range selectedUtxos {
		if isConfidential {
			psetWithoutFees.AddBlindedInput(utxo.Hash(), utxo.Index(), &partial.ConfidentialWitnessUtxo{
				AssetCommitment: utxo.AssetCommitment(),
				ValueCommitment: utxo.ValueCommitment(),
				Script:          utxo.Script(),
				Nonce:           utxo.Nonce(),
				RangeProof:      utxo.RangeProof(),
				SurjectionProof: utxo.SurjectionProof(),
			}, nil)
			blindingPrivKeysOfInputs = append(blindingPrivKeysOfInputs, bk)
		} else {
			psetWithoutFees.AddInput(utxo.Hash(), utxo.Index(), &partial.WitnessUtxo{Asset: asset, Value: utxo.Value(), Script: fromScript}, nil)
		}
	}
	psetWithoutFees.AddOutput(asset, uint64(amount), toScript, false)
	psetWithoutFees.AddOutput(asset, uint64(change), fromScript, false)

	b64, err := psetWithoutFees.Data.ToBase64()
	if err != nil {
		return fmt.Errorf("base64: %w", err)
	}

	client := taxi.NewTaxiProtobufClient(taxiURL, &http.Client{})
	replyEstimate, err := client.GetAssetEstimate(context.Background(), &taxi.GetAssetEstimateRequest{
		Unsigned: &taxi.Unsigned{
			Pset:       b64,
			SatPerByte: 0.1,
		},
		AssetHash: asset,
	})
	if err != nil {
		return fmt.Errorf("taxi: %w", err)
	}
	if replyEstimate.AssetAmount <= 0 {
		return fmt.Errorf("taxi: Invalid number")
	}

	// Delete last inserted element that in our case is the change
	psetWithoutFees.Data.Outputs = psetWithoutFees.Data.Outputs[:1]
	psetWithoutFees.Data.UnsignedTx.Outputs = psetWithoutFees.Data.UnsignedTx.Outputs[:1]

	// Add recalculated change after estimate with taxi
	psetWithoutFees.AddOutput(asset, uint64(change)-replyEstimate.AssetAmount, fromScript, false)
	// Add dummy Fee output
	psetWithoutFees.AddOutput(network.Regtest.AssetID, 0, []byte{}, false)

	// Get all blinding keys in place
	toBlindKey, err := confidential.ToBlindingKey(to, currentNetwork)
	if err != nil {
		return fmt.Errorf("toBlindKey: %w", err)
	}
	fromBlindKeyPair, err := keypair.FromPrivateKey(blinding)
	if err != nil {
		return fmt.Errorf("fromBlindKeyPair: %w", err)
	}
	blindingPubKeysOfOutputs := [][]byte{toBlindKey, fromBlindKeyPair.PublicKey.SerializeCompressed()}

	// Blind
	err = psetWithoutFees.BlindWithKeys(blindingPrivKeysOfInputs, blindingPubKeysOfOutputs)
	if err != nil {
		return fmt.Errorf("Blinding: %w", err)
	}

	// Delete last inserted element that in this case is the dummy fee output
	psetWithoutFees.Data.Outputs = psetWithoutFees.Data.Outputs[:2]
	psetWithoutFees.Data.UnsignedTx.Outputs = psetWithoutFees.Data.UnsignedTx.Outputs[:2]

	b64, err = psetWithoutFees.Data.ToBase64()
	if err != nil {
		return fmt.Errorf("base64: %w", err)
	}

	fmt.Println()
	fmt.Println(b64)

	return nil
}
