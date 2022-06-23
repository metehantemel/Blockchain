package blockchain

import "bytes"

type TransactionInput struct {
	TransactionID []byte `bson:"ti_id"`
	Vout          int    `bson:"vout"`
	Signature     []byte `bson:"signature"`
	PublicKey     []byte `bson:"public_key"`
}

func (_input *TransactionInput) UsesKey(_publicKeyHash []byte) bool {
	_lockingHash := HashPublicKey(_input.PublicKey)

	return bytes.Compare(_lockingHash, _publicKeyHash) == 0
}
