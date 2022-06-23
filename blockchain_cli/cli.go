package blockchain_cli

import (
	"flag"
	"fmt"
	"os"
)

type CLI struct{}

func (_cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  createwallet - Generates a new key-pair and saves it into the wallet file")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  listaddresses - Lists all addresses from the wallet file")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")
}

func (_cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		_cli.printUsage()
		os.Exit(1)
	}
}

func (_cli *CLI) Run() {
	_cli.validateArgs()

	_getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	_printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	_sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	_createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	_createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	_listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)

	_getBalanceAddress := _getBalanceCmd.String("address", "", "The address to get balance for")
	_createBlockchainAddress := _createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	_sendFrom := _sendCmd.String("from", "", "Source wallet address")
	_sendTo := _sendCmd.String("to", "", "Destination wallet address")
	_sendAmount := _sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		_error := _getBalanceCmd.Parse(os.Args[2:])
		if _error != nil {
			panic(_error)
		}
	case "printchain":
		_error := _printChainCmd.Parse(os.Args[2:])
		if _error != nil {
			panic(_error)
		}
	case "send":
		_error := _sendCmd.Parse(os.Args[2:])
		if _error != nil {
			panic(_error)
		}
	case "createblockchain":
		_error := _createBlockchainCmd.Parse(os.Args[2:])
		if _error != nil {
			panic(_error)
		}
	case "createwallet":
		_error := _createWalletCmd.Parse(os.Args[2:])
		if _error != nil {
			panic(_error)
		}
	case "listaddresses":
		_error := _listAddressesCmd.Parse(os.Args[2:])
		if _error != nil {
			panic(_error)
		}
	default:
		_cli.printUsage()
		os.Exit(1)
	}

	if _printChainCmd.Parsed() {
		_cli.printChain()
	}

	if _getBalanceCmd.Parsed() {
		if *_getBalanceAddress == "" {
			_getBalanceCmd.Usage()
			os.Exit(1)
		}
		_cli.getBalance(*_getBalanceAddress)
	}

	if _createBlockchainCmd.Parsed() {
		if *_createBlockchainAddress == "" {
			_createBlockchainCmd.Usage()
			os.Exit(1)
		}
		_cli.createBlockchain(*_createBlockchainAddress)
	}

	if _sendCmd.Parsed() {
		if *_sendFrom == "" || *_sendTo == "" || *_sendAmount <= 0 {
			_sendCmd.Usage()
			os.Exit(1)
		}

		_cli.send(*_sendFrom, *_sendTo, *_sendAmount)
	}

	if _createWalletCmd.Parsed() {
		_cli.createWallet()
	}

	if _listAddressesCmd.Parsed() {
		_cli.listAddresses()
	}
}
