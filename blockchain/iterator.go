package blockchain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlockchainIterator struct {
	currentHash []byte
	cursor      *mongo.Cursor
}

func (_blockchain *Blockchain) Iterator() *BlockchainIterator {
	_options := options.Find().SetSort(bson.M{"_id": -1})
	_cursor, _error := _blockchain.DB.collection.Find(context.TODO(), bson.D{}, _options)
	if _error != nil {
		panic(_error)
	}

	_iterator := &BlockchainIterator{_blockchain.LastBlockHash, _cursor}

	return _iterator
}

func (_iterator *BlockchainIterator) Next() *Block {
	var _block *Block

	_iterator.cursor.Next(context.TODO())

	_error := _iterator.cursor.Decode(&_block)
	if _error != nil {
		panic(_error)
	}

	return _block
}
