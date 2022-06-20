package blockchain

type Blockchain struct {
	LastBlockHash []byte
	DB            *db
}

func (_blockchain *Blockchain) AddBlock(_data string) {
	_previousBlock := _blockchain.DB.GetLastBlock()
	newBlock := NewBlock(_data, _previousBlock.Hash)
	_blockchain.DB.AddBlock(newBlock)
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func NewBlockChain() *Blockchain {
	_db := Create_DB()
	var _lastBlockHash []byte

	if _db.GetBlockCount() == 0 {
		_genesisBlock := NewGenesisBlock()
		_lastBlockHash = _genesisBlock.Hash
		_db.AddBlock(_genesisBlock)
	} else {
		_lastBlock := _db.GetLastBlock()
		_lastBlockHash = _lastBlock.Hash
	}

	return &Blockchain{
		_lastBlockHash,
		_db,
	}
}
