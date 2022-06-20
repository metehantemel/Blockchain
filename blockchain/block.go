package blockchain

import (
	"time"
)

type Block struct {
	Timestamp         int64  `bson:"timestamp"`
	Data              []byte `bson:"data"`
	PreviousBlockHash []byte `bson:"previous_hash"`
	Hash              []byte `bson:"hash"`
	Nonce             int    `bson:"nonce"`
}

func NewBlock(_data string, _previousBlockHash []byte) *Block {
	_block := &Block{
		time.Now().Unix(),
		[]byte(_data),
		_previousBlockHash,
		[]byte{},
		0,
	}

	_proofOfWork := NewProofOfWork(_block)
	_nonce, _hash := _proofOfWork.Run()

	_block.Hash = _hash[:]
	_block.Nonce = _nonce

	return _block
}
