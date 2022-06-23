package blockchain

import "bytes"

type TransactionOutput struct {
	Value         int    `bson:"to_value"`
	PublicKeyHash []byte `bson:"public_key_hash"`
}

func (_out *TransactionOutput) Lock(_address []byte) {
	_publicKeyHash := Base58Decode(_address)
	_publicKeyHash = _publicKeyHash[1 : len(_publicKeyHash)-4]
	_out.PublicKeyHash = _publicKeyHash
}

func (_out *TransactionOutput) IsLockedWithKey(_publicKeyHash []byte) bool {
	return bytes.Compare(_out.PublicKeyHash, _publicKeyHash) == 0
}

func NewTransactionOutput(_value int, _address string) *TransactionOutput {
	_transactionOutput := &TransactionOutput{
		_value,
		nil,
	}
	_transactionOutput.Lock([]byte(_address))

	return _transactionOutput
}
