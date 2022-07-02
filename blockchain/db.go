package blockchain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MONGO_URI = "mongodb://docker:mongopw@localhost:49156"
const DB_BLOCKCHAIN = "blockchain"
const COLLECTION_BLOCKS = "blocks"
const COLLECTION_UTXOS = "chainstate"

type db struct {
	client          *mongo.Client
	collection      *mongo.Collection
	utxo_collection *mongo.Collection
}

func Create_DB() *db {
	_mongoClient, _error := mongo.Connect(context.TODO(), options.Client().ApplyURI(MONGO_URI))
	if _error != nil {
		panic(_error)
	}
	_db := &db{
		_mongoClient,
		_mongoClient.Database(DB_BLOCKCHAIN).Collection(COLLECTION_BLOCKS),
		_mongoClient.Database(DB_BLOCKCHAIN).Collection(COLLECTION_UTXOS),
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
	_db.utxo_collection = _db.client.Database(DB_BLOCKCHAIN).Collection(COLLECTION_UTXOS)
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
	return _blocks
}

func (_db db) GetLastBlock() *Block {
	var _block Block
	_options := options.FindOne().SetSort(bson.M{"$natural": -1})
	_db.collection.FindOne(context.TODO(), bson.D{}, _options).Decode(&_block)
	return &_block
}

func (_db db) DropUTXOs() bool {
	_error := _db.utxo_collection.Drop(context.TODO())
	if _error != nil {
		panic(_error)
		return false
	}
	return true
}

func (_db db) AddUTXO(_UTXO *UTXO) {
	_result, _error := _db.utxo_collection.InsertOne(
		context.TODO(),
		_UTXO,
	)
	if _error != nil {
		panic(_error)
	}
	_ = _result
}

func (_db db) GetUTXO(_transactionID []byte) UTXO {
	var _UTXO UTXO
	_db.utxo_collection.FindOne(context.TODO(), bson.M{"t_id": _transactionID}).Decode(&_UTXO)
	return _UTXO
}

func (_db db) DeleteUTXO(_transactionID []byte) {
	_result, _error := _db.utxo_collection.DeleteOne(context.TODO(), bson.M{"t_id": _transactionID})
	_ = _result
	_ = _error
}

func (_db db) PutUTXO(_transactionID []byte, _transactionOutputs []TransactionOutput) {
	_filter := bson.D{{"t_id", _transactionID}}
	_update := bson.D{{"$set", bson.D{{"transaction_outputs", _transactionOutputs}}}}
	_count, _error := _db.utxo_collection.CountDocuments(context.TODO(), _filter)
	_ = _error
	if _count == 0 {
		_UTXO := UTXO{
			_transactionID,
			_transactionOutputs,
		}
		_db.AddUTXO(&_UTXO)
	} else {
		_result, _error := _db.utxo_collection.UpdateOne(context.TODO(), _filter, _update)
		_ = _result
		_ = _error
	}
}

func (_db db) GetUTXOCount() int64 {
	_count, _error := _db.utxo_collection.EstimatedDocumentCount(context.TODO())
	if _error != nil {
		panic(_error)
	}
	return _count
}
