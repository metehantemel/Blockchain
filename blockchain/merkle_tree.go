package blockchain

import "crypto/sha256"

type MerkleTree struct {
	RootNode *MerkleNode
}

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

func NewMerkleNode(_left, _right *MerkleNode, _data []byte) *MerkleNode {
	_merkleNode := MerkleNode{}

	if _left == nil && _right == nil {
		_hash := sha256.Sum256(_data)
		_merkleNode.Data = _hash[:]
	} else {
		_previousHashes := append(_left.Data, _right.Data...)
		_hash := sha256.Sum256(_previousHashes)
		_merkleNode.Data = _hash[:]
	}

	_merkleNode.Left = _left
	_merkleNode.Right = _right

	return &_merkleNode
}

func NewMerkleTree(_data [][]byte) *MerkleTree {
	var _nodes []MerkleNode

	if len(_data)%2 != 0 {
		_data = append(_data, _data[len(_data)-1])
	}

	for _, _datum := range _data {
		_node := NewMerkleNode(nil, nil, _datum)
		_nodes = append(_nodes, *_node)
	}

	for _i := 0; _i < len(_data)/2; _i++ {
		var _newLevel []MerkleNode

		for _j := 0; _j < len(_nodes); _j += 2 {
			_node := NewMerkleNode(&_nodes[_j], &_nodes[_j+1], nil)
			_newLevel = append(_newLevel, *_node)
		}

		_nodes = _newLevel
	}

	_merkleTree := MerkleTree{&_nodes[0]}

	return &_merkleTree
}
