package main

import (
	"log"

	"github.com/leaanthony/clir"
)

func main() {
	// Declaration of flag variables
	var from, to, asset, psetBase64Topup, psetBase64Sign, privateKey string
	var amount int
	var regtest bool = false
	var explorerURL string = "http://localhost:3001"
	var taxiURL string = "http://localhost:8080"

	// Create new cli
	cli := clir.NewCli("Taxi CLI", "A command line interface to work with Liquid Taxi", "v0.0.1")

	// Create subcommand
	createCmd := cli.NewSubCommand("create", "Create an unsigned PSET without LBTC fees")

	// Required flags
	createCmd.StringFlag("from", "(REQUIRED) From address used to retrive utxos and send back the change", &from)
	createCmd.StringFlag("to", "(REQUIRED) Recipient address to send coins", &to)
	createCmd.StringFlag("asset", "(REQUIRED) Asset hash to include in the transaction", &asset)
	createCmd.IntFlag("amount", "(REQUIRED) Amount of coins to send to recipient", &amount)
	// Optional Flags
	createCmd.BoolFlag("regtest", "Work with local regtest", &regtest)
	createCmd.StringFlag("taxi", "Taxi endpoint API", &taxiURL)
	createCmd.StringFlag("explorer", "Electrs REST API", &explorerURL)
	// Action
	createCmd.Action(func() error {
		return createAction(from, to, asset, amount, regtest, explorerURL, taxiURL)
	})

	// Topup subcommand
	topupCommand := cli.NewSubCommand("topup", "Topup given PSET with LBTC fees")
	// Required flags
	topupCommand.StringFlag("pset", "(REQUIRED) Partial Signed Elements Transaction base64 encoded", &psetBase64Topup)
	topupCommand.StringFlag("asset", "(REQUIRED) Asset hash to include in the transaction", &asset)
	// Optional Flags
	topupCommand.StringFlag("taxi", "Taxi endpoint API", &taxiURL)
	// Action
	topupCommand.Action(func() error {
		return topupAction(psetBase64Topup, asset, taxiURL)
	})

	// Sign subcommand
	signCommand := cli.NewSubCommand("sign", "Sign given PSET with private key (hex format)")
	// Required flags
	signCommand.StringFlag("pset", "(REQUIRED) Partial Signed Elements Transaction base64 encoded", &psetBase64Sign)
	signCommand.StringFlag("key", "(REQUIRED) EC Private Key (hex encoded)", &privateKey)
	// Optional Flags
	signCommand.BoolFlag("regtest", "Work with local regtest", &regtest)
	// Action
	signCommand.Action(func() error {
		return signAction(psetBase64Sign, privateKey, regtest)
	})

	// Run the application
	err := cli.Run()
	if err != nil {
		log.Fatal(err)
	}
}
