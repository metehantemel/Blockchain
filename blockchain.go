package main

type Blockchain struct {
	blocks []*Block
}

func (_blockchain *Blockchain) AddBlock(data string) {
	previousBlock := _blockchain.blocks[len(_blockchain.blocks)-1]
	newBlock := NewBlock(data, previousBlock.Hash)
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
