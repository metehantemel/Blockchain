package blockchain_cli

import (
	"blockchainGO/blockchain"
	"fmt"
	"log"
)

func (_cli *CLI) createBlockchain(_address string) {
	if !blockchain.ValidateAddress(_address) {
		log.Panic("ERROR: Address is not valid")
	}

	_blockchain := blockchain.NewBlockChain(_address)
	_ = _blockchain
	fmt.Println("Done!")
}
