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

	if _db.GetBlockCount() == 0 {
		_coinbaseTransaction := NewCoinbaseTransaction(_address, "genesisCoinbaseData")
		_genesisBlock := NewGenesisBlock(_coinbaseTransaction)
		_lastBlockHash = _genesisBlock.Hash
		_db.AddBlock(_genesisBlock)
	} else {
		_lastBlock := _db.GetLastBlock()
		_lastBlockHash = _lastBlock.Hash
	}

	return &Blockchain{
		_lastBlockHash,
		_db,
	}
}

func (_blockchain *Blockchain) MineBlock(_transactions []*Transaction) {
	_previousBlock := _blockchain.DB.GetLastBlock()

	for _, _transaction := range _transactions {
		if _blockchain.VerifyTransaction(_transaction) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	newBlock := NewBlock(_transactions, _previousBlock.Hash)
	_blockchain.DB.AddBlock(newBlock)
}

func (_blockchain *Blockchain) FindSpendableOutputs(_publicKeyHash []byte, _amount int) (int, map[string][]int) {
	_unspentOutputs := make(map[string][]int)
	_unspentTransactions := _blockchain.FindUnspentTransactions(_publicKeyHash)
	_accumulated := 0

Work:
	for _, _transaction := range _unspentTransactions {
		_transactionID := hex.EncodeToString(_transaction.ID)

		for _outputID, _output := range _transaction.Vout {
			if _output.IsLockedWithKey(_publicKeyHash) && _accumulated < _amount {
				_accumulated += _output.Value
				_unspentOutputs[_transactionID] = append(_unspentOutputs[_transactionID], _outputID)

				if _accumulated >= _amount {
					break Work
				}
			}
		}
	}

	return _accumulated, _unspentOutputs
}

func (_blockchain *Blockchain) FindUnspentTransactions(_publicKeyHash []byte) []Transaction {
	var _unspentTransactions []Transaction
	_spentTransactionOutputs := make(map[string][]int)
	_blockchainIterator := _blockchain.Iterator()

	for {
		_block := _blockchainIterator.Next()

		for _, _transcation := range _block.Transactions {
			_transcationID := hex.EncodeToString(_transcation.ID)

		Outputs:
			for _outputID, _output := range _transcation.Vout {
				if _spentTransactionOutputs != nil {
					for _, _spentOutputID := range _spentTransactionOutputs[_transcationID] {
						if _spentOutputID == _outputID {
							continue Outputs
						}
					}
				}

				if _output.IsLockedWithKey(_publicKeyHash) {
					_unspentTransactions = append(_unspentTransactions, *_transcation)
				}
			}

			if _transcation.IsCoinbase() == false {
				for _, _transactionInput := range _transcation.Vin {
					if _transactionInput.UsesKey(_publicKeyHash) {
						_transactionInputID := hex.EncodeToString(_transactionInput.TransactionID)
						_spentTransactionOutputs[_transactionInputID] = append(_spentTransactionOutputs[_transactionInputID], _transactionInput.Vout)
					}
				}
			}
		}

		if len(_block.PreviousBlockHash) == 0 {
			break
		}
	}

	return _unspentTransactions
}

func (_blockchain *Blockchain) FindUTXO(_publicKeyHash []byte) []TransactionOutput {
	var _UTXOs []TransactionOutput
	_unspentTransactions := _blockchain.FindUnspentTransactions(_publicKeyHash)

	for _, _transaction := range _unspentTransactions {
		for _, _output := range _transaction.Vout {
			if _output.IsLockedWithKey(_publicKeyHash) {
				_UTXOs = append(_UTXOs, _output)
			}
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
