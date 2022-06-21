package blockchain

import "encoding/hex"

type Blockchain struct {
	LastBlockHash []byte
	DB            *db
}

func (_blockchain *Blockchain) MineBlock(_transactions []*Transaction) {
	_previousBlock := _blockchain.DB.GetLastBlock()
	newBlock := NewBlock(_transactions, _previousBlock.Hash)
	_blockchain.DB.AddBlock(newBlock)
}

func NewGenesisBlock(_coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{_coinbase}, []byte{})
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

func (_blockchain *Blockchain) FindUnspentTransactions(_address string) []Transaction {
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
					for _, _spentOutput := range _spentTransactionOutputs[_transcationID] {
						if _spentOutput == _outputID {
							continue Outputs
						}
					}
				}

				if _output.CanBeUnlockedWith(_address) {
					_unspentTransactions = append(_unspentTransactions, *_transcation)
				}
			}

			if _transcation.IsCoinbase() == false {
				for _, _transactionInput := range _transcation.Vin {
					if _transactionInput.CanUnlockOutputWith(_address) {
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

func (_blockchain *Blockchain) FindUTXO(_address string) []TransactionOutput {
	var _UTXOs []TransactionOutput
	_unspentTransactions := _blockchain.FindUnspentTransactions(_address)

	for _, _transaction := range _unspentTransactions {
		for _, _output := range _transaction.Vout {
			if _output.CanBeUnlockedWith(_address) {
				_UTXOs = append(_UTXOs, _output)
			}
		}
	}

	return _UTXOs
}
