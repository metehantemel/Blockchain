package main

import (
	"time"
)

type Block struct {
	Timestamp         int64
	Data              []byte
	PreviousBlockHash []byte
	Hash              []byte
	Nonce             int
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
