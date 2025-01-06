package main

import (
	"fmt"
	"log"
)

func (cli *CLI) Listaddresses() {
	wallets, err := NewWallets()

	if err != nil {
		log.Panic(err)
	}

	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}
