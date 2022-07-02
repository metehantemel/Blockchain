package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"log"
)

type Blockchain struct {
	LastBlockHash []byte
	DB            *db
}

func NewBlockChain(_address string) *Blockchain {
	_db := Create_DB()
	var _lastBlockHash []byte
	var _blockchain Blockchain
	if _db.GetBlockCount() == 0 {
		_coinbaseTransaction := NewCoinbaseTransaction(_address, "genesisCoinbaseData")
		_genesisBlock := NewGenesisBlock(_coinbaseTransaction)
		_lastBlockHash = _genesisBlock.Hash
		_db.AddBlock(_genesisBlock)
		_blockchain = Blockchain{
			_lastBlockHash,
			_db,
		}
		_UTXOSet := UTXOSet{&_blockchain}
		_UTXOSet.ReIndex()
	} else {
		_lastBlock := _db.GetLastBlock()
		_lastBlockHash = _lastBlock.Hash
		_blockchain = Blockchain{
			_lastBlockHash,
			_db,
		}
	}

	return &_blockchain
}

func (_blockchain *Blockchain) MineBlock(_transactions []*Transaction) *Block {
	_previousBlock := _blockchain.DB.GetLastBlock()

	for _, _transaction := range _transactions {
		if _blockchain.VerifyTransaction(_transaction) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	newBlock := NewBlock(_transactions, _previousBlock.Hash)
	_blockchain.DB.AddBlock(newBlock)
	return newBlock
}

func (_blockchain *Blockchain) FindUTXOs() []UTXO {
	_UTXOs := []UTXO{}
	_spentTransactionOuts := make(map[string][]int)
	_iterator := _blockchain.Iterator()

	for {
		_block := _iterator.Next()
		for _, _transaction := range _block.Transactions {
			_transactionID := hex.EncodeToString(_transaction.ID)

			_UTXOc := UTXO{[]byte{}, []TransactionOutput{}}
			_UTXOs = append(_UTXOs, _UTXOc)
			_UTXO := &_UTXOs[len(_UTXOs)-1]
			_UTXO.TransactionID = _transaction.ID

		Outputs:
			for _outId, _out := range _transaction.Vout {
				if _spentTransactionOuts[_transactionID] != nil {
					for _, _spentTransactionOutId := range _spentTransactionOuts[_transactionID] {
						if _spentTransactionOutId == _outId {
							continue Outputs
						}
					}
				}

				_UTXO.Outputs = append(_UTXO.Outputs, _out)

			}

			if _transaction.IsCoinbase() == false {
				for _, _input := range _transaction.Vin {
					_inputID := hex.EncodeToString(_input.TransactionID)
					_spentTransactionOuts[_inputID] = append(_spentTransactionOuts[_inputID], _input.Vout)
				}
			}

		}

		if len(_block.PreviousBlockHash) == 0 {
			break
		}
	}

	return _UTXOs
}

func (_blockchain *Blockchain) FindTransaction(_ID []byte) (Transaction, error) {
	_bcIterator := _blockchain.Iterator()
	for {
		_block := _bcIterator.Next()

		for _, _transaction := range _block.Transactions {
			if bytes.Compare(_transaction.ID, _ID) == 0 {
				return *_transaction, nil
			}
		}

		if len(_block.PreviousBlockHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction not found")
}

func (_blockchain *Blockchain) SignTransaction(_transaction *Transaction, _privateKey ecdsa.PrivateKey) {
	_previousTransactions := make(map[string]Transaction)

	for _, _vin := range _transaction.Vin {
		_previousTransaction, _error := _blockchain.FindTransaction(_vin.TransactionID)
		if _error != nil {
			panic(_error)
		}
		_previousTransactions[hex.EncodeToString(_previousTransaction.ID)] = _previousTransaction
	}

	_transaction.Sign(_privateKey, _previousTransactions)
}

func (_blockchain *Blockchain) VerifyTransaction(_transaction *Transaction) bool {
	if _transaction.IsCoinbase() {
		return true
	}

	_previousTransactions := make(map[string]Transaction)

	for _, _vin := range _transaction.Vin {
		_previousTransaction, _error := _blockchain.FindTransaction(_vin.TransactionID)
		if _error != nil {
			panic(_error)
		}
		_previousTransactions[hex.EncodeToString(_previousTransaction.ID)] = _previousTransaction
	}

	return _transaction.Verify(_previousTransactions)
}
