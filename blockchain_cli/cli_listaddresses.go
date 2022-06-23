package blockchain_cli

import (
	"blockchainGO/blockchain"
	"fmt"
	"log"
)

func (_cli *CLI) listAddresses() {
	_wallets, _error := blockchain.NewWallets()
	if _error != nil {
		log.Panic(_error)
	}

	_addresses := _wallets.GetAddresses()

	for _, _address := range _addresses {
		_wallet := _wallets.GetWallet(_address)
		_balance := _wallet.GetBalance()
		fmt.Printf("----------\r\nWallet: %s \r\nBalance: %d\r\n", _address, _balance)
	}
}
