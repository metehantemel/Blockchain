package blockchain_cli

import (
	"blockchainGO/blockchain"
	"fmt"
)

func (_cli *CLI) reindexUTXO() {
	_blockchain := blockchain.NewBlockChain("")
	_UTXOSet := blockchain.UTXOSet{_blockchain}
	_UTXOSet.ReIndex()

	_count := _UTXOSet.CountTransactions()
	fmt.Printf("There are %d transactions in UTXO set.\r\n", _count)
}
