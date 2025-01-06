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

	bc := NewBlockChain()
	UTXOSet := UTXOSet{bc}

	defer bc.db.Close()

	tx := NewUTXOTransaction(from, to, amount, &UTXOSet)
	cbTx := NewCoinbaseTx(from, "")
	newBlock := bc.MineBlock([]*Transaction{cbTx, tx})

	UTXOSet.Update(newBlock)

	fmt.Println("Success !")
}
