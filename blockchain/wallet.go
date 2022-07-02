package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const walletFile = "wallet.dat"
const addressChecksumLen = 4

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	_curve := elliptic.P256()
	_private, _error := ecdsa.GenerateKey(_curve, rand.Reader)
	if _error != nil {
		panic(_error)
	}
	_publicKey := append(_private.PublicKey.X.Bytes(), _private.PublicKey.Y.Bytes()...)

	return *_private, _publicKey
}

func NewWallet() *Wallet {
	_private, _public := newKeyPair()
	_wallet := Wallet{_private, _public}

	return &_wallet
}

func (_wallet Wallet) GetAddress() []byte {
	_publicKeyHash := HashPublicKey(_wallet.PublicKey)

	_versionedPayload := append([]byte{version}, _publicKeyHash...)
	_checksum := checksum(_versionedPayload)

	_fullPayload := append(_versionedPayload, _checksum...)
	_address := Base58Encode(_fullPayload)

	return _address
}

func HashPublicKey(_publicKey []byte) []byte {
	_publicSHA256 := sha256.Sum256(_publicKey)

	_RIPEMD160Hasher := ripemd160.New()
	_, _error := _RIPEMD160Hasher.Write(_publicSHA256[:])
	if _error != nil {
		panic(_error)
	}
	_publicRIPEMD160 := _RIPEMD160Hasher.Sum(nil)

	return _publicRIPEMD160
}

func checksum(_payload []byte) []byte {
	_firstSHA := sha256.Sum256(_payload)
	_secondSHA := sha256.Sum256(_firstSHA[:])

	return _secondSHA[:addressChecksumLen]
}

func ValidateAddress(_address string) bool {
	_publicKeyHash := Base58Decode([]byte(_address))
	_actualChecksum := _publicKeyHash[len(_publicKeyHash)-addressChecksumLen:]
	_version := _publicKeyHash[0]
	_publicKeyHash = _publicKeyHash[1 : len(_publicKeyHash)-addressChecksumLen]
	_targetChecksum := checksum(append([]byte{_version}, _publicKeyHash...))

	return bytes.Compare(_actualChecksum, _targetChecksum) == 0
}

func (_wallet *Wallet) GetBalance() int {
	_address := string(_wallet.GetAddress())
	_blockchain := NewBlockChain(_address)
	_UTXOSet := UTXOSet{_blockchain}

	_balance := 0
	_publicKeyHash := Base58Decode([]byte(_address))
	_publicKeyHash = _publicKeyHash[1 : len(_publicKeyHash)-4]
	_UTXOs := _UTXOSet.FindUTXO(_publicKeyHash)
	fmt.Printf("%d Length UTXOS\r\n", len(_UTXOs))

	for _, _out := range _UTXOs {
		_balance += _out.Value
	}

	return _balance
}
