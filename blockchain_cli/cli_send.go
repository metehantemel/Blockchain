package blockchain_cli

import (
	"blockchainGO/blockchain"
	"fmt"
	"log"
)

func (_cli *CLI) send(_from, _to string, _amount int) {
	if !blockchain.ValidateAddress(_from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !blockchain.ValidateAddress(_to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	_blockchain := blockchain.NewBlockChain(_from)
	_UTXOSet := blockchain.UTXOSet{_blockchain}

	_transaction := blockchain.NewUTXOTransaction(_from, _to, _amount, &_UTXOSet)
	_coinbaseTransaction := blockchain.NewCoinbaseTransaction(_from, "")
	_transactions := []*blockchain.Transaction{
		_coinbaseTransaction,
		_transaction,
	}

	_block := _blockchain.MineBlock(_transactions)
	_UTXOSet.Update(_block)
	fmt.Println("Success!")
}
