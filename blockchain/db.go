package blockchain

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MONGO_URI = "mongodb://docker:mongopw@localhost:49153"
const DB_BLOCKCHAIN = "blockchain"
const COLLECTION_BLOCKS = "blocks"

type db struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func Create_DB() *db {
	_mongoClient, _error := mongo.Connect(context.TODO(), options.Client().ApplyURI(MONGO_URI))
	if _error != nil {
		panic(_error)
	}
	_db := &db{
		_mongoClient,
		_mongoClient.Database(DB_BLOCKCHAIN).Collection(COLLECTION_BLOCKS),
	}
	return _db
}

func (_db db) connect() {
	_mongoClient, _error := mongo.Connect(context.TODO(), options.Client().ApplyURI(MONGO_URI))
	if _error != nil {
		panic(_error)
	}
	_db.client = _mongoClient
	_db.collection = _db.client.Database(DB_BLOCKCHAIN).Collection(COLLECTION_BLOCKS)
}

func (_db db) disconnect() {
	if _db.client != nil {
		_db.client.Disconnect(context.TODO())
	}
}

func (_db db) AddBlock(_block *Block) {
	_result, _error := _db.collection.InsertOne(
		context.TODO(),
		_block,
	)
	if _error != nil {
		panic(_error)
	}
	_ = _result
}

func (_db db) GetBlockCount() int64 {
	_count, _error := _db.collection.EstimatedDocumentCount(context.TODO())
	if _error != nil {
		panic(_error)
	}
	return _count
}

func (_db db) GetBlocks() []*Block {
	_blocks := []*Block{}
	_cursor, _error := _db.collection.Find(context.TODO(), bson.D{})
	if _error != nil {
		panic(_error)
	}

	_error = _cursor.All(context.TODO(), &_blocks)
	if _error != nil {
		panic(_error)
	}
	fmt.Printf("%d", len(_blocks))
	return _blocks
}

func (_db db) GetLastBlock() *Block {
	var _block Block
	_options := options.FindOne().SetSort(bson.M{"$natural": -1})
	_db.collection.FindOne(context.TODO(), bson.D{}, _options).Decode(&_block)
	return &_block
}
