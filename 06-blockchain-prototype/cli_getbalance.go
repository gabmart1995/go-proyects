package main

import (
	"fmt"
	"log"
)

func (cli *CLI) getBalance(address string, nodeID string) {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}

	bc := NewBlockChain(nodeID)
	defer bc.db.Close()

	UTXOSet := UTXOSet{bc}

	balance := 0
	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]

	UTXOs := UTXOSet.FindUTXO(pubKeyHash)
	// fmt.Println(UTXOs)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
