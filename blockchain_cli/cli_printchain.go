package blockchain_cli

import (
	"blockchainGO/blockchain"
	"fmt"
	"strconv"
)

func (_cli *CLI) printChain() {

	_iterator := blockchain.NewBlockChain("").Iterator()

	for {
		_block := _iterator.Next()

		fmt.Printf("====== Block %x ======\n", _block.Hash)
		fmt.Printf("Prev. hash: %x\n", _block.PreviousBlockHash)
		fmt.Printf("Hash: %x\n", _block.Hash)
		pow := blockchain.NewProofOfWork(_block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		for _, _transaction := range _block.Transactions {
			fmt.Println(_transaction)
		}
		fmt.Println()

		if len(_block.PreviousBlockHash) == 0 {
			break
		}
	}
}
