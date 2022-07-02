package blockchain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (_utxoSet UTXOSet) Iterator() *mongo.Cursor {
	_options := options.Find().SetSort(bson.M{"_id": -1})

	_cursor, _error := _utxoSet.Blockchain.DB.utxo_collection.Find(context.TODO(), bson.D{}, _options)
	if _error != nil {
		panic(_error)
	}

	return _cursor
}
