package main

import (
	"fmt"
	"log"
)

func (cli *CLI) Listaddresses(nodeID string) {
	wallets, err := NewWallets(nodeID)

	if err != nil {
		log.Panic(err)
	}

	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}
