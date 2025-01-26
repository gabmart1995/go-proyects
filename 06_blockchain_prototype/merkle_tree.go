package main

import "crypto/sha256"

/**
Modulo especializado en utilizar un algoritmo de
optimizacion de bloques en el blockchain. Implementado
por Satoshi Nakamoto en bitcoin
*/

type MerkleTree struct {
	RootNode *MerkleNode
}

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode

	// si el nodo es impar la transaccion se duplica
	if (len(data) % 2) != 0 {
		data = append(data, data[len(data)-1])
	}

	// inicializamos las hojas
	for _, datum := range data {
		node := NewMerkleNode(nil, nil, datum)
		nodes = append(nodes, *node)
	}

	for i := 0; i < (len(data) / 2); i++ {
		var newLevel []MerkleNode

		// extraemos los nodos pares
		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, *node)
		}

		nodes = newLevel
	}

	mTree := MerkleTree{RootNode: &nodes[0]}

	return &mTree
}

// construye un node merkle
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	mNode := MerkleNode{}
	var hash [32]byte

	// si no existe nodos en los laterales es una hoja
	// se genera un hash de transaccion con la data
	// sino se toma la informacion de los nodos laterales
	if left == nil && right == nil {
		hash = sha256.Sum256(data)

	} else {
		prevHashes := append(left.Data, right.Data...)
		hash = sha256.Sum256(prevHashes)

	}

	mNode.Data = hash[:]
	mNode.Left = left
	mNode.Right = right

	return &mNode
}
