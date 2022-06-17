package main

type Blockchain struct {
	blocks []*Block
}

func (_blockchain *Blockchain) AddBlock(_data string) {
	_previousBlock := _blockchain.blocks[len(_blockchain.blocks)-1]
	newBlock := NewBlock(_data, _previousBlock.Hash)
	_blockchain.blocks = append(_blockchain.blocks, newBlock)
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func NewBlockChain() *Blockchain {
	return &Blockchain{
		[]*Block{NewGenesisBlock()},
	}
}
