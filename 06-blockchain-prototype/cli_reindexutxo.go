package main

import "fmt"

func (cli *CLI) ReindexUTXO(nodeID string) {
	bc := NewBlockChain(nodeID)
	UTXOSet := UTXOSet{bc}
	UTXOSet.Reindex()

	count := UTXOSet.CountTransactions()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
}
