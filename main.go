package main

import (
	"blockchainGO/blockchain"
)

func main() {
	_blockchain := blockchain.NewBlockChain()

	_cli := CLI{_blockchain}
	_cli.Run()
}
