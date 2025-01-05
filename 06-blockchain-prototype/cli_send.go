package main

import (
	"fmt"
	"log"
)

func (cli *CLI) send(from, to string, amount int) {
	if !ValidateAddress(from) {
		log.Panic("ERROR: sender address is not valid")
	}

	if !ValidateAddress(to) {
		log.Panic("ERROR: recipent address is not valid")
	}

	bc := NewBlockChain(from)
	UTXOSet := UTXOSet{bc}

	defer bc.db.Close()

	tx := NewUTXOTransaction(from, to, amount, &UTXOSet)
	newBlock := bc.MineBlock([]*Transaction{tx})

	UTXOSet.Update(newBlock)

	fmt.Println("Success !")
}
