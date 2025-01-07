package main

import (
	"bytes"
	"encoding/gob"
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

func StartServer(nodeID, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	miningAddress = minerAddress

	// abrimos el puerto usando el prtocolo http
	// cada nodo especifica un puerto
	ln, err := net.Listen(protocol, nodeAddress)
	defer ln.Close()

	if err != nil {
		log.Panic(err)
	}

	bc := NewBlockChain(nodeID)

	if nodeAddress != knownNodes[0] {
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
		if b == 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
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
	defer conn.Close()

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

func requestBlocks() {
	for _, node := range knownNodes {
		sendGetBlocks(node)
	}
}

func sendGetBlocks(address string) {
	payload := gobEncode(getblocks{AddrFrom: nodeAddress})
	request := append(commandToBytes("getblocks"), payload...)

	sendData(address, request)
}

func sendGetData(address, kind string, id []byte) {
	payload := gobEncode(getdata{nodeAddress, kind, id})
	request := append(commandToBytes("getdata"), payload...)

	sendData(address, request)
}
