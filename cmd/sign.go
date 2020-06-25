package main

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcutil"
	"github.com/tiero/ocean/pkg/keypair"
	"github.com/tiero/ocean/pkg/partial"
	"github.com/vulpemventures/go-elements/network"
	"github.com/vulpemventures/go-elements/payment"
	"github.com/vulpemventures/go-elements/pset"
)

func signAction(psetBase64 string, privateKey string, regtest bool) error {
	if psetBase64 == "" || privateKey == "" {
		return fmt.Errorf("missing required flag")
	}
	// Network
	currentNetwork := network.Liquid
	if regtest {
		currentNetwork = network.Regtest
	}

	decoded, err := pset.NewPsetFromBase64(psetBase64)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	var keyPair *keypair.KeyPair
	keyPair, err = keypair.FromPrivateKey(privateKey)
	if err != nil {
		// Try to see if it's WIF
		wif, err2 := btcutil.DecodeWIF(privateKey)
		if err2 != nil {
			return fmt.Errorf("private key: %w", err)
		}
		keyPair, _ = keypair.FromPrivateKey(hex.EncodeToString(wif.PrivKey.Serialize()))
	}
	pay := payment.FromPublicKey(keyPair.PublicKey, &currentNetwork, nil)

	psetWithFees := &partial.Partial{Data: decoded, Network: &currentNetwork}
	for i := 0; i < len(psetWithFees.Data.Inputs); i++ {
		currInput := psetWithFees.Data.Inputs[i]
		if bytes.Equal(currInput.WitnessUtxo.Script, pay.Script) {
			err := psetWithFees.SignWithPrivateKey(i, keyPair)
			if err != nil {
				return fmt.Errorf("SignWithPrivateKey: %w", err)
			}
		}
	}
	pFinalized := psetWithFees.Data
	err = pset.FinalizeAll(pFinalized)
	if err != nil {
		return fmt.Errorf("FinalizeAll: %w", err)
	}

	if !pFinalized.IsComplete() {
		return fmt.Errorf("pset not complete: %w", err)
	}

	err = pFinalized.SanityCheck()
	if err != nil {
		return fmt.Errorf("sanity check: %w", err)
	}

	b64, err := pFinalized.ToBase64()
	if err != nil {
		return fmt.Errorf("base64: %w", err)
	}

	fmt.Println(b64)

	// Extract the final signed transaction from the Pset wrapper.
	finalTx, err := pset.Extract(pFinalized)
	if err != nil {
		return fmt.Errorf("Extract: %w", err)
	}

	// Serialize the transaction and try to broadcast.
	txHex, err := finalTx.ToHex()
	if err != nil {
		return fmt.Errorf("ToHex: %w", err)
	}

	fmt.Println()
	fmt.Println(txHex)

	return nil
}
