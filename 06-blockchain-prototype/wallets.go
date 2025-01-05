package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"log"
	"os"
)

// almacena una coleccion de wallets
type Wallets struct {
	Wallets map[string]*Wallet
}

func NewWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFromFile()

	return &wallets, err
}

// carga los wallets desde el archivo
func (ws *Wallets) LoadFromFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := os.ReadFile(walletFile)

	if err != nil {
		log.Panic(err)
	}

	var wallets Wallets

	/*if err := json.Unmarshal(fileContent, &wallets); err != nil {
		log.Panic(err)
	}*/

	// OLD CODE go 1.18
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	if err := decoder.Decode(&wallets); err != nil {
		log.Panic(err)
	}

	ws.Wallets = wallets.Wallets

	return nil
}

// GetWallet returns a Wallet by its address
func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

// crea un wallet y lo asigna a los wallets
func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := string(wallet.GetAddress())

	ws.Wallets[address] = wallet

	return address
}

func (ws Wallets) SaveToFile() {
	//jsonData, err := json.Marshal(ws)

	// OLD CODE go 1.18
	var content bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)

	if err != nil {
		log.Panic(err)
	}

	err = os.WriteFile(walletFile, content.Bytes(), 0644)
	// err = os.WriteFile(walletFile, jsonData, 0644)

	if err != nil {
		log.Panic(err)
	}
}
