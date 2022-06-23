package blockchain

import (
	"bytes"
	"crypto/sha256"
	"time"
)

type Block struct {
	Timestamp         int64          `bson:"timestamp"`
	Transactions      []*Transaction `bson:"transactions"`
	PreviousBlockHash []byte         `bson:"previous_hash"`
	Hash              []byte         `bson:"hash"`
	Nonce             int            `bson:"nonce"`
}

func NewBlock(_transactions []*Transaction, _previousBlockHash []byte) *Block {
	_block := &Block{
		time.Now().Unix(),
		_transactions,
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

func NewGenesisBlock(_coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{_coinbase}, []byte{})
}

func (_block *Block) HashTransactions() []byte {
	var _transactionHashes [][]byte
	var _transactionHash [32]byte

	for _, _transaction := range _block.Transactions {
		_transactionHashes = append(_transactionHashes, _transaction.ID)
	}
	_transactionHash = sha256.Sum256(bytes.Join(_transactionHashes, []byte{}))

	return _transactionHash[:]
}
