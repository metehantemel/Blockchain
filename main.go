package main

import (
	"fmt"
	"strconv"
)

func main() {
	_blockchain := NewBlockChain()

	_blockchain.AddBlock("Send 1BTC To Test")
	_blockchain.AddBlock("Send 2BTC To Test")

	for _, _block := range _blockchain.blocks {
		fmt.Printf("Previous hash: %x\n", _block.PreviousBlockHash)
		fmt.Printf("Timestamp: %d\n", _block.Timestamp)
		fmt.Printf("Data: %s\n", _block.Data)
		fmt.Printf("Hash: %x\n", _block.Hash)
		_proofOfWork := NewProofOfWork(_block)
		fmt.Printf("Valid: %s\n", strconv.FormatBool(_proofOfWork.Validate()))
		fmt.Printf("\n")
	}
}
