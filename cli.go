package main

import (
	"blockchainGO/blockchain"
	"flag"
	"fmt"
	"os"
	"strconv"
)

type CLI struct {
	blockchain *blockchain.Blockchain
}

func (_cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  addblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("  printchain - print all the blocks of the blockchain")
}

func (_cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		_cli.printUsage()
		os.Exit(1)
	}
}

func (_cli *CLI) Run() {
	_cli.validateArgs()

	_addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	_printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	_addBlockData := _addBlockCmd.String("data", "", "Block data")

	switch os.Args[1] {
	case "addblock":
		_error := _addBlockCmd.Parse(os.Args[2:])
		if _error != nil {
			panic(_error)
		}
	case "printchain":
		_error := _printChainCmd.Parse(os.Args[2:])
		if _error != nil {
			panic(_error)
		}
	default:
		_cli.printUsage()
		os.Exit(1)
	}

	if _addBlockCmd.Parsed() {
		if *_addBlockData == "" {
			_addBlockCmd.Usage()
			os.Exit(1)
		}
		_cli.addBlock(*_addBlockData)
	}

	if _printChainCmd.Parsed() {
		_cli.printChain()
	}
}

func (_cli *CLI) addBlock(_data string) {
	_cli.blockchain.AddBlock(_data)
	fmt.Println("Success!")
}

func (_cli *CLI) printChain() {

	_iterator := _cli.blockchain.Iterator()

	for {
		_block := _iterator.Next()

		fmt.Printf("Prev. hash: %x\n", _block.PreviousBlockHash)
		fmt.Printf("Data: %s\n", _block.Data)
		fmt.Printf("Hash: %x\n", _block.Hash)
		pow := blockchain.NewProofOfWork(_block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(_block.PreviousBlockHash) == 0 {
			break
		}
	}
}
