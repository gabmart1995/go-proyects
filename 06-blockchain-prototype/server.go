package main

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
)

/*
* modulo de creacion de nodos de red

Son servidores especializados que manejan la propia
informacion de la red de blockachain desentralizada.
*/
var nodeAddress string
var knownNodes = []string{"localhost:3000"}
var miningAddress string
var blocksInTransit = [][]byte{}
var mempool = make(map[string]Transaction)

const (
	protocol      = "tcp"
	nodeVersion   = 1
	commandLength = 12
)

type verzion struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

type addr struct {
	AddrList []string
}

type block struct {
	AddrFrom string
	Block    []byte
}

type getblocks struct {
	AddrFrom string
}

type getdata struct {
	AddrFrom string
	Type     string
	ID       []byte
}

type tx struct {
	AddFrom     string
	Transaction []byte
}

type inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

func StartServer(nodeID, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	miningAddress = minerAddress

	// abrimos el puerto usando el prtocolo http
	// cada nodo especifica un puerto
	ln, err := net.Listen(protocol, nodeAddress)

	if err != nil {
		log.Panic(err)
	}

	defer ln.Close()

	bc := NewBlockChain(nodeID)

	if nodeAddress != knownNodes[0] {
		// se manda peticion al nodo raiz para actualizar la
		// version
		sendVersion(knownNodes[0], bc)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Panic(err)
		}

		go handleConnection(conn, bc)
	}
}

// prepara la peticion de la version
func sendVersion(addr string, bc *BlockChain) {
	bestHeight := bc.GetBestHeight()
	payload := gobEncode(verzion{
		Version:    nodeVersion,
		BestHeight: bestHeight,
		AddrFrom:   nodeAddress,
	})

	request := append(commandToBytes("version"), payload...)

	// fmt.Println(string(commandToBytes("version")))

	sendData(addr, request)
}

// parsea el comando a binario
func commandToBytes(command string) []byte {
	var bytes [commandLength]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return string(command)
}

// tranforma la informacion recibida por tcp
func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	encoder := gob.NewEncoder(&buff)

	if err := encoder.Encode(data); err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// abre la conexion y manda los datos
func sendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)

	if err != nil {
		fmt.Printf("%s is not available\n", addr)

		// si falla la conexion debemos actualizar los nodos
		// de conexion
		var updatedNodes []string

		for _, node := range knownNodes {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}

		knownNodes = updatedNodes

		return
	}

	defer conn.Close()

	// copiamos los bytes al resto de la red
	_, err = io.Copy(conn, bytes.NewReader(data))

	if err != nil {
		log.Panic(err)
	}

}

// controlador de peticiones tcp
func handleConnection(conn net.Conn, bc *BlockChain) {
	request, err := io.ReadAll(conn)

	if err != nil {
		log.Panic(err)
	}

	command := bytesToCommand(request[:commandLength])
	fmt.Printf("Received %s command\n", command)

	switch command {
	case "addr":
		handleAddr(request)

	case "block":
		handleBlock(request, bc)

	case "inv":
		handleInv(request)

	case "getblocks":
		handleGetBlocks(request, bc)

	case "getdata":
		handleGetData(request, bc)

	case "tx":
		handleTx(request, bc)

	case "version":
		handleVersion(request, bc)

	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
}

func handleAddr(request []byte) {
	var (
		buff    bytes.Buffer
		payload addr
	)

	_, err := buff.Write(request[:commandLength])

	if err != nil {
		log.Panic(err)
	}

	decoder := gob.NewDecoder(&buff)

	if err := decoder.Decode(&payload); err != nil {
		log.Panic(err)
	}

	knownNodes = append(knownNodes, payload.AddrList...)
	fmt.Printf("There are %d known nodes now!\n", len(knownNodes))

	requestBlocks()
}

func handleBlock(request []byte, bc *BlockChain) {
	var (
		buff    bytes.Buffer
		payload block
	)

	_, err := buff.Write(request[:commandLength])

	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Block
	block := DeserializeBlock(blockData)

	// en este punto se recibe el bloque correctamente
	fmt.Println("Recevied a new block!")
	bc.AddBlock(block)

	fmt.Printf("Added block %x\n", block.Hash)

	// arma la peticion en el buffer
	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		blocksInTransit = blocksInTransit[1:]

	} else {
		UTXOSet := UTXOSet{Blockchain: bc}
		UTXOSet.Reindex()

	}
}

func handleInv(request []byte) {
	var (
		buff    bytes.Buffer
		payload inv
	)

	if _, err := buff.Write(request[commandLength:]); err != nil {
		log.Panic(err)
	}

	decoder := gob.NewDecoder(&buff)
	if err := decoder.Decode(&payload); err != nil {
		log.Panic(err)
	}

	fmt.Printf("Recevied inventory with %d %s\n", len(payload.Items), payload.Type)

	if payload.Type == "block" {
		blocksInTransit = payload.Items
		blockHash := payload.Items[0]

		sendGetData(payload.AddrFrom, "block", blockHash)

		newInTransit := [][]byte{}
		for _, b := range blocksInTransit {
			if !bytes.Equal(b, blockHash) {
				newInTransit = append(newInTransit, b)
			}
		}

		blocksInTransit = newInTransit
	}

	if payload.Type == "tx" {
		txID := payload.Items[0]

		if mempool[hex.EncodeToString(txID)].ID == nil {
			sendGetData(payload.AddrFrom, "tx", txID)
		}
	}
}

func handleGetBlocks(request []byte, bc *BlockChain) {
	var (
		buff    bytes.Buffer
		payload getblocks
	)

	if _, err := buff.Write(request[commandLength:]); err != nil {
		log.Panic(err)
	}

	decoder := gob.NewDecoder(&buff)
	if err := decoder.Decode(&payload); err != nil {
		log.Panic(err)
	}

	blocks := bc.GetBlockHashes()
	sendInv(payload.AddrFrom, "block", blocks)
}

func handleGetData(request []byte, bc *BlockChain) {
	var (
		buff    bytes.Buffer
		payload getdata
	)

	if _, err := buff.Write(request[commandLength:]); err != nil {
		log.Panic(err)
	}

	decoder := gob.NewDecoder(&buff)
	if err := decoder.Decode(&payload); err != nil {
		log.Panic(err)
	}

	if payload.Type == "block" {
		block, err := bc.GetBlock([]byte(payload.ID))

		if err != nil {
			return
		}

		sendBlock(payload.AddrFrom, &block)
	}

	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := mempool[txID]

		sendTx(payload.AddrFrom, &tx)
	}
}

func handleTx(request []byte, bc *BlockChain) {
	var (
		buff    bytes.Buffer
		payload tx
	)

	if _, err := buff.Write(request[commandLength:]); err != nil {
		log.Panic(err)
	}

	decoder := gob.NewDecoder(&buff)
	if err := decoder.Decode(&payload); err != nil {
		log.Panic(err)
	}

	txData := payload.Transaction
	tx := DeserializeTransaction(txData)

	mempool[hex.EncodeToString(tx.ID)] = tx

	if nodeAddress == knownNodes[0] {
		for _, node := range knownNodes {
			if node != nodeAddress && node != payload.AddFrom {
				sendInv(node, "tx", [][]byte{tx.ID})
			}
		}

	} else {
		if len(mempool) >= 2 && len(miningAddress) > 0 {
		MineTransactions:

			var txs []*Transaction

			for id := range mempool {
				tx := mempool[id]

				if bc.VerifyTransaction(&tx) {
					txs = append(txs, &tx)
				}
			}

			if len(txs) == 0 {
				fmt.Println("All transactions are invalid! Waiting for new ones...")
				return
			}

			cbTx := NewCoinbaseTx(miningAddress, "")
			txs = append(txs, cbTx)

			newBlock := bc.MineBlock(txs)
			UTXOSet := UTXOSet{
				Blockchain: bc,
			}
			UTXOSet.Reindex()

			fmt.Println("New block is mined!")

			for _, tx := range txs {
				txID := hex.EncodeToString(tx.ID)
				delete(mempool, txID)
			}

			for _, node := range knownNodes {
				if node != nodeAddress {
					sendInv(node, "block", [][]byte{newBlock.Hash})
				}
			}

			if len(mempool) > 0 {
				goto MineTransactions
			}
		}
	}
}

func handleVersion(request []byte, bc *BlockChain) {
	var (
		buff    bytes.Buffer
		payload verzion
	)

	if _, err := buff.Write(request[commandLength:]); err != nil {
		log.Panic(err)
	}

	decoder := gob.NewDecoder(&buff)
	if err := decoder.Decode(&payload); err != nil {
		log.Panic(err)
	}

	myBestHeight := bc.GetBestHeight()
	foreignerBestHeight := payload.BestHeight

	if myBestHeight < foreignerBestHeight {
		sendGetBlocks(payload.AddrFrom)

	} else if myBestHeight > foreignerBestHeight {
		sendVersion(payload.AddrFrom, bc)

	}

	if !nodeIsKnown(payload.AddrFrom) {
		knownNodes = append(knownNodes, payload.AddrFrom)
	}
}

func requestBlocks() {
	for _, node := range knownNodes {
		sendGetBlocks(node)
	}
}

func nodeIsKnown(addr string) bool {
	for _, node := range knownNodes {
		if node == addr {
			return true
		}
	}

	return false
}

func sendGetBlocks(address string) {
	payload := gobEncode(getblocks{AddrFrom: nodeAddress})
	request := append(commandToBytes("getblocks"), payload...)

	sendData(address, request)
}

func sendGetData(address, kind string, id []byte) {
	payload := gobEncode(getdata{
		AddrFrom: nodeAddress,
		Type:     kind,
		ID:       id,
	})
	request := append(commandToBytes("getdata"), payload...)

	sendData(address, request)
}

func sendTx(addr string, tnx *Transaction) {
	data := tx{
		AddFrom:     nodeAddress,
		Transaction: tnx.Serialize(),
	}
	payload := gobEncode(data)
	request := append(commandToBytes("tx"), payload...)

	sendData(addr, request)
}

func sendInv(address, kind string, items [][]byte) {
	inventory := inv{
		AddrFrom: nodeAddress,
		Type:     kind,
		Items:    items,
	}
	payload := gobEncode(inventory)
	request := append(commandToBytes("inv"), payload...)

	sendData(address, request)
}

func sendBlock(addr string, b *Block) {
	data := block{
		AddrFrom: nodeAddress,
		Block:    b.Serialize(),
	}

	payload := gobEncode(data)
	request := append(commandToBytes("block"), payload...)

	sendData(addr, request)
}
