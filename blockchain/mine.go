package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

const difficulty = 16

type ProofOfWork struct {
	block      *Block
	difficulty *big.Int
}

func NewProofOfWork(_block *Block) *ProofOfWork {
	_difficulty := big.NewInt(1)
	_difficulty.Lsh(
		_difficulty,
		uint(256-difficulty))

	_proofOfWork := &ProofOfWork{
		_block,
		_difficulty,
	}

	return _proofOfWork
}

func (_proofOfWork *ProofOfWork) prepareData(_nonce int) []byte {
	_data := bytes.Join(
		[][]byte{
			_proofOfWork.block.PreviousBlockHash,
			_proofOfWork.block.HashTransactions(),
			IntToHex(_proofOfWork.block.Timestamp),
			IntToHex(int64(difficulty)),
			IntToHex(int64(_nonce)),
		},
		[]byte{},
	)

	return _data
}

func (_proofOfWork *ProofOfWork) Run() (int, []byte) {
	var _hashInt big.Int
	var _hash [32]byte
	_maxNonce := math.MaxInt64
	_nonce := 0

	fmt.Printf("Mining block \"%x\"\n", _proofOfWork.block.HashTransactions())

	for _nonce < _maxNonce {
		_data := _proofOfWork.prepareData(_nonce)
		_hash = sha256.Sum256(_data)
		_hashInt.SetBytes(_hash[:])

		if _hashInt.Cmp(_proofOfWork.difficulty) == -1 {
			fmt.Printf("\rSolved: %x", _hash)
			break
		} else {
			_nonce++
		}
	}
	fmt.Print("\n\n")
	return _nonce, _hash[:]
}

func (_proofOfWork *ProofOfWork) Validate() bool {
	var _hashInt big.Int

	_data := _proofOfWork.prepareData(_proofOfWork.block.Nonce)
	_hash := sha256.Sum256(_data)
	_hashInt.SetBytes(_hash[:])

	_isValid := _hashInt.Cmp(_proofOfWork.difficulty) == -1

	return _isValid
}
