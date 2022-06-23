package blockchain

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
)

type Wallets struct {
	Wallets map[string]*Wallet
}

func NewWallets() (*Wallets, error) {
	_wallets := Wallets{}
	_wallets.Wallets = make(map[string]*Wallet)

	_error := _wallets.LoadFromFile()

	return &_wallets, _error
}

func (_wallets *Wallets) CreateWallet() string {
	_wallet := NewWallet()
	_address := fmt.Sprintf("%s", _wallet.GetAddress())

	_wallets.Wallets[_address] = _wallet

	return _address
}

func (_wallets *Wallets) GetAddresses() []string {
	var _addresses []string

	for _address := range _wallets.Wallets {
		_addresses = append(_addresses, _address)
	}

	return _addresses
}

func (_wallets Wallets) GetWallet(_address string) Wallet {
	return *_wallets.Wallets[_address]
}

func (_wallets *Wallets) LoadFromFile() error {
	if _, _error := os.Stat(walletFile); os.IsNotExist(_error) {
		return _error
	}

	_fileContent, _error := ioutil.ReadFile(walletFile)
	if _error != nil {
		panic(_error)
	}

	var _ws Wallets
	gob.Register(elliptic.P256())
	_decoder := gob.NewDecoder(bytes.NewReader(_fileContent))
	_error = _decoder.Decode(&_ws)
	if _error != nil {
		panic(_error)
	}

	_wallets.Wallets = _ws.Wallets

	return nil
}

func (_wallets Wallets) SaveToFile() {
	var _content bytes.Buffer

	gob.Register(elliptic.P256())

	_encoder := gob.NewEncoder(&_content)
	_error := _encoder.Encode(_wallets)
	if _error != nil {
		panic(_error)
	}

	_error = ioutil.WriteFile(walletFile, _content.Bytes(), 0644)
	if _error != nil {
		panic(_error)
	}
}
