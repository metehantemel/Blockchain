package main

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	Timestamp         int64
	Data              []byte
	PreviousBlockHash []byte
	Hash              []byte
}

func (_block *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(_block.Timestamp, 10))
	headers := bytes.Join(
		[][]byte{_block.PreviousBlockHash, _block.Data, timestamp},
		[]byte{},
	)
	hash := sha256.Sum256(headers)

	_block.Hash = hash[:]
}

func NewBlock(data string, previousBlockHash []byte) *Block {
	block := &Block{
		time.Now().Unix(),
		[]byte(data),
		previousBlockHash,
		[]byte{}}

	block.SetHash()

	return block
}
