package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain_%s.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type BlockChain struct {
	tip []byte
	db  *bolt.DB
}

func (bc *BlockChain) MineBlock(transactions []*Transaction) *Block {
	var lastHash []byte

	for _, tx := range transactions {
		if !bc.VerifyTransaction(tx) {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	// obtiene el ultimo bloque
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		lastHash = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	newBlock := NewBlock(transactions, lastHash)

	// se ultiliza el ultimo hash para crear el nuevo bloque
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if err := b.Put(newBlock.Hash, newBlock.Serialize()); err != nil {
			log.Fatal(err)
		}

		if err := b.Put([]byte("l"), newBlock.Hash); err != nil {
			log.Fatal(err)
		}

		bc.tip = newBlock.Hash

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return newBlock
}

func NewBlockChain(nodeID string) *BlockChain {
	dbFile := fmt.Sprintf(dbFile, nodeID)

	if !dbExists(dbFile) {
		fmt.Println("No existing blockchain found. Create one first")
		os.Exit(1)
	}

	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := BlockChain{tip, db}

	return &bc
}

func CreateBlockchain(address, nodeID string) *BlockChain {
	dbFile := fmt.Sprintf(dbFile, nodeID)

	if dbExists(dbFile) {
		fmt.Println("Blockchain alredy exists")
		os.Exit(1)
	}

	cbtx := NewCoinbaseTx(address, genesisCoinbaseData)
	genesis := NewGenesisBlock(cbtx)

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		if err := b.Put(genesis.Hash, genesis.Serialize()); err != nil {
			log.Panic(err)
		}

		if err := b.Put([]byte("l"), genesis.Hash); err != nil {
			log.Panic(err)
		}

		tip = genesis.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := BlockChain{tip, db}

	return &bc
}

// carga los bloques en memoria y se ejecutan cuando se necesiten
func (bc *BlockChain) Iterator() *BlockChainIterator {
	bci := &BlockChainIterator{
		currentHash: bc.tip,
		db:          bc.db,
	}

	return bci
}

func dbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

// retorna todas las salidas de transacciones que no fueron gastadas
func (bc *BlockChain) FindUTXO() map[string]TXOutputs {
	UTXO := make(map[string]TXOutputs)
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// verificamos si la salda fue gastada
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				// establece la transaccion
				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}

			if !tx.IsCoinbase() {
				for _, in := range tx.Vin {
					inTXID := hex.EncodeToString(in.Txid)
					spentTXOs[inTXID] = append(spentTXOs[inTXID], in.Vout)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return UTXO
}

// agrupa las transacciones y por ultimo las firma
func (bc *BlockChain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)

		if err != nil {
			log.Panic(err)
		}

		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	// firma la transaccion
	tx.Sign(privKey, prevTXs)
}

// Localiza la transaccion por ID
func (bc *BlockChain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Equal(tx.ID, ID) {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction is not found")
}

func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)

		if err != nil {
			log.Panic(err)
		}

		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}

// retorna la altura del ultimo bloque
func (bc *BlockChain) GetBestHeight() int {
	var lastBlock Block

	err := bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		lastHash := bucket.Get([]byte("l"))
		blockData := bucket.Get(lastHash)

		lastBlock = *DeserializeBlock(blockData)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return lastBlock.Height
}

func (bc *BlockChain) AddBlock(block *Block) {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		blockInDb := bucket.Get(block.Hash)

		if blockInDb != nil {
			return nil
		}

		// en caso de no estar se serializa nuevamente
		blockData := block.Serialize()

		if err := bucket.Put(block.Hash, blockData); err != nil {
			log.Panic(err)
		}

		lastHash := bucket.Get([]byte("l"))
		lastBlockData := bucket.Get(lastHash)
		lastBlock := DeserializeBlock(lastBlockData)

		if block.Height > lastBlock.Height {
			if err := bucket.Put([]byte("l"), block.Hash); err != nil {
				log.Panic(err)
			}

			bc.tip = block.Hash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

// obtiene los hashes de cada bloque
func (bc *BlockChain) GetBlockHashes() [][]byte {
	var blocks [][]byte

	bci := bc.Iterator()

	for {
		block := bci.Next()
		blocks = append(blocks, block.Hash)

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return blocks
}

func (bc *BlockChain) GetBlock(blockHash []byte) (Block, error) {
	var block Block

	err := bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		blocksData := bucket.Get(blockHash)

		if blocksData == nil {
			return errors.New("Block is not found")
		}

		block = *DeserializeBlock(blocksData)

		return nil
	})

	if err != nil {
		return block, err
	}

	return block, nil
}
