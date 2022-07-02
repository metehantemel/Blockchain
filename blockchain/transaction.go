package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
)

const subsidy = 10

type Transaction struct {
	ID   []byte              `bson:"t_id"`
	Vin  []TransactionInput  `bson:"transaction_input"`
	Vout []TransactionOutput `bson:"transaction_output"`
}

func NewCoinbaseTransaction(_to, _data string) *Transaction {
	if _data == "" {
		_randData := make([]byte, 20)
		_, _error := rand.Read(_randData)
		if _error != nil {
			log.Panic(_error)
		}
		_data = fmt.Sprintf("%x", _randData)
		fmt.Printf("DATA : %s\r\n", _data)
	}

	_transactionInput := TransactionInput{
		[]byte{},
		-1,
		nil,
		[]byte(_data),
	}
	_transactionOutput := NewTransactionOutput(subsidy, _to)
	_transaction := Transaction{
		nil,
		[]TransactionInput{_transactionInput},
		[]TransactionOutput{*_transactionOutput},
	}
	_transaction.ID = _transaction.Hash()

	return &_transaction
}

func NewUTXOTransaction(_from, _to string, _amount int, _UTXOSet *UTXOSet) *Transaction {
	var _inputs []TransactionInput
	var _outputs []TransactionOutput

	_wallets, _error := NewWallets()
	if _error != nil {
		panic(_error)
	}
	_wallet := _wallets.GetWallet(_from)
	_publicKeyHash := HashPublicKey(_wallet.PublicKey)
	_account, _validOutputs := _UTXOSet.FindSpendableOutputs(_publicKeyHash, _amount)

	if _account < _amount {
		log.Panic("ERROR: Not enough funds")
	}

	for _transactionID, _outs := range _validOutputs {
		_transactionID, _error := hex.DecodeString(_transactionID)
		if _error != nil {
			panic(_error)
		}

		for _, _out := range _outs {
			_input := TransactionInput{_transactionID, _out, nil, _wallet.PublicKey}
			_inputs = append(_inputs, _input)
		}
	}

	_outputs = append(_outputs, *NewTransactionOutput(_amount, _to))
	if _account > _amount {
		_outputs = append(_outputs, *NewTransactionOutput(_account-_amount, _from))
	}

	_transaction := Transaction{nil, _inputs, _outputs}
	_transaction.ID = _transaction.Hash()
	_UTXOSet.Blockchain.SignTransaction(&_transaction, _wallet.PrivateKey)

	return &_transaction
}

func (_transaction *Transaction) IsCoinbase() bool {
	return len(_transaction.Vin) == 1 && len(_transaction.Vin[0].TransactionID) == 0 && _transaction.Vin[0].Vout == -1
}

func (_transaction Transaction) Serialize() []byte {
	var _encoded bytes.Buffer

	_encoder := gob.NewEncoder(&_encoded)
	_error := _encoder.Encode(_transaction)
	if _error != nil {
		log.Panic(_error)
	}

	return _encoded.Bytes()
}

func (_transaction *Transaction) Hash() []byte {
	var _hash [32]byte

	_transactionCopy := *_transaction
	_transactionCopy.ID = []byte{}

	_hash = sha256.Sum256(_transactionCopy.Serialize())

	return _hash[:]
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

func (_transaction *Transaction) TrimmedCopy() Transaction {
	var _inputs []TransactionInput
	var _outputs []TransactionOutput

	for _, _vin := range _transaction.Vin {
		_inputs = append(
			_inputs,
			TransactionInput{
				_vin.TransactionID,
				_vin.Vout,
				nil,
				nil},
		)
	}

	for _, _vout := range _transaction.Vout {
		_outputs = append(
			_outputs,
			TransactionOutput{
				_vout.Value,
				_vout.PublicKeyHash},
		)
	}

	_transactionCopy := Transaction{
		_transaction.ID,
		_inputs,
		_outputs}

	return _transactionCopy
}

func (_transaction *Transaction) Sign(_privateKey ecdsa.PrivateKey, _previousTransactions map[string]Transaction) {
	if _transaction.IsCoinbase() {
		return
	}

	_transactionCopy := _transaction.TrimmedCopy()

	for _inputID, _vin := range _transactionCopy.Vin {
		_previousTransaction := _previousTransactions[hex.EncodeToString(_vin.TransactionID)]
		_transactionCopy.Vin[_inputID].Signature = nil
		_transactionCopy.Vin[_inputID].PublicKey = _previousTransaction.Vout[_vin.Vout].PublicKeyHash
		_transactionCopy.ID = _transactionCopy.Hash()
		_transactionCopy.Vin[_inputID].PublicKey = nil

		_r, _s, _error := ecdsa.Sign(rand.Reader, &_privateKey, _transactionCopy.ID)
		if _error != nil {
			panic(_error)
		}
		_signature := append(_r.Bytes(), _s.Bytes()...)

		_transaction.Vin[_inputID].Signature = _signature
	}
}

func (_transaction *Transaction) Verify(_previousTransactions map[string]Transaction) bool {
	_transactionCopy := _transaction.TrimmedCopy()
	_curve := elliptic.P256()

	for _inputID, _vin := range _transaction.Vin {
		_previousTransaction := _previousTransactions[hex.EncodeToString(_vin.TransactionID)]
		_transactionCopy.Vin[_inputID].Signature = nil
		_transactionCopy.Vin[_inputID].PublicKey = _previousTransaction.Vout[_vin.Vout].PublicKeyHash
		_transactionCopy.ID = _transactionCopy.Hash()
		_transactionCopy.Vin[_inputID].PublicKey = nil

		_r := big.Int{}
		_s := big.Int{}
		_signatureLen := len(_vin.Signature)
		_r.SetBytes(_vin.Signature[:(_signatureLen / 2)])
		_s.SetBytes(_vin.Signature[(_signatureLen / 2):])

		_x := big.Int{}
		_y := big.Int{}
		_keyLen := len(_vin.PublicKey)
		_x.SetBytes(_vin.PublicKey[:(_keyLen / 2)])
		_y.SetBytes(_vin.PublicKey[(_keyLen / 2):])

		_rawPublicKey := ecdsa.PublicKey{_curve, &_x, &_y}
		if ecdsa.Verify(&_rawPublicKey, _transactionCopy.ID, &_r, &_s) == false {
			return false
		}
	}

	return true
}

func (_transaction Transaction) String() string {
	var _lines []string

	_lines = append(_lines, fmt.Sprintf("--- Transaction %x:", _transaction.ID))

	for i, _input := range _transaction.Vin {

		_lines = append(_lines, fmt.Sprintf("     Input %d:", i))
		_lines = append(_lines, fmt.Sprintf("       TXID:      %x", _input.TransactionID))
		_lines = append(_lines, fmt.Sprintf("       Out:       %d", _input.Vout))
		_lines = append(_lines, fmt.Sprintf("       Signature: %x", _input.Signature))
		_lines = append(_lines, fmt.Sprintf("       PubKey:    %x", _input.PublicKey))
	}

	for _i, _output := range _transaction.Vout {
		_lines = append(_lines, fmt.Sprintf("     Output %d:", _i))
		_lines = append(_lines, fmt.Sprintf("       Value:  %d", _output.Value))
		_lines = append(_lines, fmt.Sprintf("       Script: %x", _output.PublicKeyHash))
	}

	return strings.Join(_lines, "\n")
}
