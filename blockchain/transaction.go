package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

const subsidy = 10

type Transaction struct {
	ID   []byte              `bson:"t_id"`
	Vin  []TransactionInput  `bson:"transaction_input"`
	Vout []TransactionOutput `bson:"transaction_output"`
}

type TransactionOutput struct {
	Value        int    `bson:"to_value"`
	ScriptPubKey string `bson:"script_pub_key"`
}

type TransactionInput struct {
	TransactionID []byte `bson:"ti_id"`
	Vout          int    `bson:"vout"`
	ScriptSig     string `bson:"script_sig"`
}

func NewCoinbaseTransaction(_to, _data string) *Transaction {
	if _data == "" {
		_data = fmt.Sprintf("Reward to '%s'", _to)
	}

	_transactionInput := TransactionInput{
		[]byte{},
		-1,
		_data,
	}
	_transactionOutput := TransactionOutput{
		subsidy,
		_to,
	}
	_transaction := Transaction{
		nil,
		[]TransactionInput{_transactionInput},
		[]TransactionOutput{_transactionOutput},
	}
	_transaction.SetID()

	return &_transaction
}

func NewUTXOTransaction(_from, _to string, _amount int, _blockchain *Blockchain) *Transaction {
	var _inputs []TransactionInput
	var _outputs []TransactionOutput

	_account, _validOutputs := _blockchain.FindSpendableOutputs(_from, _amount)

	if _account < _amount {
		log.Panic("ERROR: Not enough funds")
	}

	for _transactionID, _outs := range _validOutputs {
		_transactionID, _error := hex.DecodeString(_transactionID)
		if _error != nil {
			panic(_error)
		}

		for _, _out := range _outs {
			_input := TransactionInput{_transactionID, _out, _from}
			_inputs = append(_inputs, _input)
		}
	}

	_outputs = append(_outputs, TransactionOutput{_amount, _to})
	if _account > _amount {
		_outputs = append(_outputs, TransactionOutput{_account - _amount, _from})
	}

	_transaction := Transaction{nil, _inputs, _outputs}
	_transaction.SetID()

	return &_transaction
}

func (_blockchain *Blockchain) FindSpendableOutputs(_address string, _amount int) (int, map[string][]int) {
	_unspentOutputs := make(map[string][]int)
	_unspentTransactions := _blockchain.FindUnspentTransactions(_address)
	_accumulated := 0

Work:
	for _, _transaction := range _unspentTransactions {
		_transactionID := hex.EncodeToString(_transaction.ID)

		for _outputID, _output := range _transaction.Vout {
			if _output.CanBeUnlockedWith(_address) && _accumulated < _amount {
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

func (_transaction *Transaction) IsCoinbase() bool {
	return len(_transaction.Vin) == 1 && len(_transaction.Vin[0].TransactionID) == 0 && _transaction.Vin[0].Vout == -1
}

func (_transaction *Transaction) SetID() {
	var _encoded bytes.Buffer
	var _hash [32]byte

	_encoder := gob.NewEncoder(&_encoded)
	_error := _encoder.Encode(_transaction)
	if _error != nil {
		panic(_error)
	}
	_hash = sha256.Sum256(_encoded.Bytes())
	_transaction.ID = _hash[:]
}

func (_transactionInput *TransactionInput) CanUnlockOutputWith(_unlockingData string) bool {
	return _transactionInput.ScriptSig == _unlockingData
}

func (_transactionOutput *TransactionOutput) CanBeUnlockedWith(_unlockingData string) bool {
	return _transactionOutput.ScriptPubKey == _unlockingData
}
