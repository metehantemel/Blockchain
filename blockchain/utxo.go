package blockchain

type UTXO struct {
	TransactionID []byte              `bson:"t_id"`
	Outputs       []TransactionOutput `bson:"transaction_outputs"`
}
