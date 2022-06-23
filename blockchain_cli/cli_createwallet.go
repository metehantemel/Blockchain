package blockchain_cli

import (
	"blockchainGO/blockchain"
	"fmt"
)

func (_cli *CLI) createWallet() {
	_wallets, _ := blockchain.NewWallets()
	_address := _wallets.CreateWallet()
	_wallets.SaveToFile()

	fmt.Printf("Your new address: %s\n", _address)
}
