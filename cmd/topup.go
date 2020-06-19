package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/vulpemventures/taxi-protobuf/rpc/taxi"
)

func topupAction(psetBase64 string, asset string, taxiURL string) error {
	if psetBase64 == "" || asset == "" {
		return fmt.Errorf("missing required flag")
	}
	client := taxi.NewTaxiProtobufClient(taxiURL, &http.Client{})

	replyTopup, err := client.TopupWithAsset(context.Background(), &taxi.TopupWithAssetRequest{
		Unsigned: &taxi.Unsigned{
			Pset:       psetBase64,
			SatPerByte: 0.1,
		},
		AssetHash: asset,
	})
	if err != nil {
		return fmt.Errorf("taxi: %w", err)
	}

	fmt.Println()
	fmt.Println(replyTopup.Order.Partial.Pset)

	return nil
}
