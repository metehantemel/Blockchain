package blockchain_cli

import (
	"blockchainGO/blockchain"
	"fmt"
	"log"
)

func (_cli *CLI) getBalance(_address string) {
	if !blockchain.ValidateAddress(_address) {
		log.Panic("ERROR: Address is not valid")
	}
	_wallets, _error := blockchain.NewWallets()
	if _error != nil {
		log.Panic(_error)
	}
	_walllet := _wallets.GetWallet(_address)
	_balance := _walllet.GetBalance()

	fmt.Printf("Balance of '%s': %d\n", _address, _balance)
}
