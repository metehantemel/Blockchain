package blockchain

import (
	"bytes"
	"math/big"
)

var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func Base58Encode(_input []byte) []byte {
	var _result []byte

	_x := big.NewInt(0).SetBytes(_input)

	_base := big.NewInt(int64(len(b58Alphabet)))
	_zero := big.NewInt(0)
	_mod := &big.Int{}

	for _x.Cmp(_zero) != 0 {
		_x.DivMod(_x, _base, _mod)
		_result = append(_result, b58Alphabet[_mod.Int64()])
	}

	ReverseBytes(_result)
	for _b := range _input {
		if _b == 0x00 {
			_result = append([]byte{b58Alphabet[0]}, _result...)
		} else {
			break
		}
	}

	return _result
}

func Base58Decode(_input []byte) []byte {
	_result := big.NewInt(0)
	_zeroBytes := 0

	for _b := range _input {
		if _b == 0x00 {
			_zeroBytes++
		}
	}

	_payload := _input[_zeroBytes:]
	for _, _b := range _payload {
		_charIndex := bytes.IndexByte(b58Alphabet, _b)
		_result.Mul(_result, big.NewInt(58))
		_result.Add(_result, big.NewInt(int64(_charIndex)))
	}

	_decoded := _result.Bytes()
	_decoded = append(
		bytes.Repeat([]byte{byte(0x00)}, _zeroBytes),
		_decoded...)

	return _decoded
}
