package blockchain

import (
	"bytes"
	"encoding/binary"
	"log"
)

func IntToHex(_number int64) []byte {
	_buffer := new(bytes.Buffer)
	_error := binary.Write(_buffer, binary.BigEndian, _number)

	if _error != nil {
		log.Panic(_error)
	}

	return _buffer.Bytes()
}

func ReverseBytes(_data []byte) {
	for _i, _j := 0, len(_data)-1; _i < _j; _i, _j = _i+1, _j-1 {
		_data[_i], _data[_j] = _data[_j], _data[_i]
	}
}
