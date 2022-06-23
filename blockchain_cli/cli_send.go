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

	_transaction := blockchain.NewUTXOTransaction(_from, _to, _amount, _blockchain)
	_blockchain.MineBlock([]*blockchain.Transaction{_transaction})
	fmt.Println("Success!")
}
