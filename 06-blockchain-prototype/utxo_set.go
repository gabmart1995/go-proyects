package main

import (
	"encoding/hex"
	"log"

	"github.com/boltdb/bolt"
)

const utxoBucket = "chainstate"

type UTXOSet struct {
	Blockchain *BlockChain
}

// retorna las salidas no gastadas en referencia a las entradas
func (u *UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	db := u.Blockchain.db

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(utxoBucket))
		cursor := bucket.Cursor()

		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			txID := hex.EncodeToString(key)
			outs := DeserializeOutputs(value)

			for outIdx, out := range outs.Outputs {
				isOpenTx := out.IsLockedWithKey(pubKeyHash) && accumulated < amount

				if isOpenTx {
					accumulated += out.Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return accumulated, unspentOutputs
}

// halla las transacciones usando la clave publica
func (u *UTXOSet) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	db := u.Blockchain.db

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(utxoBucket))
		cursor := bucket.Cursor()

		// bucle generador
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			outs := DeserializeOutputs(value)

			for _, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return UTXOs
}

// retorna el numero de transacciones en el set
func (u UTXOSet) CountTransactions() int {
	db := u.Blockchain.db
	counter := 0

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(utxoBucket))
		cursor := bucket.Cursor()

		// cuenta los valores
		for key, _ := cursor.First(); key != nil; key, _ = cursor.Next() {
			counter++
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return counter
}

// reconstruye el set de transacciones
func (u UTXOSet) Reindex() {
	db := u.Blockchain.db
	bucketName := []byte(utxoBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucketName)

		if err != nil && err != bolt.ErrBucketNotFound {
			log.Panic(err)
		}

		_, err = tx.CreateBucket(bucketName)

		if err != nil {
			log.Panic(err)
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	UTXO := u.Blockchain.FindUTXO()

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)

		for txId, outs := range UTXO {
			key, err := hex.DecodeString(txId)

			if err != nil {
				log.Panic(err)
			}

			err = bucket.Put(key, outs.Serialize())
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

// actualiza el conjunto de salidas de transacciones
// desde el bloque el mismo sera la propina de una cadena de bloques
func (u UTXOSet) Update(block *Block) {
	db := u.Blockchain.db

	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(utxoBucket))

		for _, tx := range block.Transactions {

			if !tx.IsCoinbase() {
				for _, vin := range tx.Vin {
					updatedOutputs := TXOutputs{}
					outsBytes := bucket.Get(vin.Txid)
					outs := DeserializeOutputs(outsBytes)

					for outIdx, out := range outs.Outputs {
						if outIdx != vin.Vout {
							updatedOutputs.Outputs = append(updatedOutputs.Outputs, out)
						}
					}

					if len(updatedOutputs.Outputs) == 0 {
						if err := bucket.Delete(vin.Txid); err != nil {
							log.Panic(err)
						}

					} else {
						if err := bucket.Put(vin.Txid, updatedOutputs.Serialize()); err != nil {
							log.Panic(err)
						}
					}
				}
			}

			// linea prototipo
			newOutputs := TXOutputs{tx.Vout}

			if err := bucket.Put(tx.ID, newOutputs.Serialize()); err != nil {
				log.Panic(err)
			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}
