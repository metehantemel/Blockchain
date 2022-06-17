package main

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
