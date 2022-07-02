package blockchain

import (
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
	var _transactions [][]byte

	for _, _transaction := range _block.Transactions {
		_transactions = append(_transactions, _transaction.Serialize())
	}
	_merkleTree := NewMerkleTree(_transactions)

	return _merkleTree.RootNode.Data
}
