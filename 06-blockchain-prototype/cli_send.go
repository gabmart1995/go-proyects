package main

import "fmt"

func (cli *CLI) send(from, to string, amount int) {
	bc := NewBlockChain(from)
	defer bc.db.Close()

	tx := NewUTXOTransaction(from, to, amount, bc)
	// fmt.Println(tx)

	bc.MineBlock([]*Transaction{tx})

	fmt.Println("Success !")
}
