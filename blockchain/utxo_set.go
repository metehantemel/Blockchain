package blockchain

import (
	"context"
	"encoding/hex"
)

type UTXOSet struct {
	Blockchain *Blockchain
}

func (_utxoSet UTXOSet) ReIndex() {
	_db := *_utxoSet.Blockchain.DB
	_db.DropUTXOs()
	_UTXOs := _utxoSet.Blockchain.FindUTXOs()
	for _, _outs := range _UTXOs {
		_db.AddUTXO(&_outs)
	}
}

func (_utxoSet UTXOSet) FindSpendableOutputs(_publicKeyHash []byte, _amount int) (int, map[string][]int) {
	_unspentOutputs := make(map[string][]int)
	_accumulated := 0

	_cursor := _utxoSet.Iterator()

	for _cursor.Next(context.TODO()) {
		var _utxo UTXO
		_error := _cursor.Decode(&_utxo)
		if _error != nil {
			panic(_error)
		}
		_transactionID := hex.EncodeToString(_utxo.TransactionID)

		for _outId, _transactionOutput := range _utxo.Outputs {
			if _transactionOutput.IsLockedWithKey(_publicKeyHash) && _accumulated < _amount {
				_accumulated += _transactionOutput.Value
				_unspentOutputs[_transactionID] = append(_unspentOutputs[_transactionID], _outId)
			}
		}
	}

	return _accumulated, _unspentOutputs
}

func (_utxoSet UTXOSet) FindUTXO(_publicKeyHash []byte) []TransactionOutput {
	var _outputs []TransactionOutput

	_cursor := _utxoSet.Iterator()

	for _cursor.Next(context.TODO()) {
		var _utxo_buffer UTXO
		_error := _cursor.Decode(&_utxo_buffer)
		if _error != nil {
			panic(_error)
		}

		for _, _transactionOutput := range _utxo_buffer.Outputs {
			if _transactionOutput.IsLockedWithKey(_publicKeyHash) {
				_outputs = append(_outputs, _transactionOutput)
			}
		}
	}

	return _outputs
}

func (_utxoSet UTXOSet) Update(_block *Block) {
	_db := _utxoSet.Blockchain.DB
	for _, _transaction := range _block.Transactions {
		if _transaction.IsCoinbase() == false {
			for _, _vin := range _transaction.Vin {
				_updatedOuts := []TransactionOutput{}
				_outs := _db.GetUTXO(_vin.TransactionID)
				for _outID, _out := range _outs.Outputs {
					if _outID != _vin.Vout {
						_updatedOuts = append(_updatedOuts, _out)
					}
				}

				if len(_updatedOuts) == 0 {
					_db.DeleteUTXO(_vin.TransactionID)
				} else {
					_db.PutUTXO(_vin.TransactionID, _updatedOuts)
				}

			}
		}

		_newOutputs := []TransactionOutput{}
		for _, _out := range _transaction.Vout {
			_newOutputs = append(_newOutputs, _out)
		}

		_db.PutUTXO(_transaction.ID, _newOutputs)
	}
}

func (_utxoSet UTXOSet) CountTransactions() int64 {
	return _utxoSet.Blockchain.DB.GetUTXOCount()
}
