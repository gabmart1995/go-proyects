package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
)

const subsidy = 10

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)

	if err := enc.Encode(tx); err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

// hashea la transaccion
func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

func NewCoinbaseTx(to, data string) *Transaction {
	if len(data) == 0 {
		randData := make([]byte, 20)

		if _, err := rand.Read(randData); err != nil {
			log.Panic(err)
		}

		data = fmt.Sprintf("%x", randData)
	}

	txin := TXInput{
		Txid:      []byte{},
		Vout:      -1,
		Signature: nil,
		PubKey:    []byte(data),
	}

	txout := NewTXOuput(subsidy, to)

	tx := Transaction{
		ID:   nil,
		Vin:  []TXInput{txin},
		Vout: []TXOutput{*txout},
	}

	tx.ID = tx.Hash()

	return &tx
}

// verifica si la transaccion es en moneda base
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 &&
		len(tx.Vin[0].Txid) == 0 &&
		tx.Vin[0].Vout == -1
}

// crea una nueva transaction
func NewUTXOTransaction(wallet *Wallet, to string, amount int, UTXOSet *UTXOSet) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	pubKeyHash := HashPubKey(wallet.PublicKey)
	acc, validOutputs := UTXOSet.FindSpendableOutputs(pubKeyHash, amount)

	// verificamos saldo
	if acc < amount {
		log.Panic("Error: Not enough founds")
	}

	// construimos las entradas
	for txId, outs := range validOutputs {
		txID, err := hex.DecodeString(txId)

		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{
				Txid:      txID,
				Vout:      out,
				Signature: nil,
				PubKey:    wallet.PublicKey,
			}

			inputs = append(inputs, input)
		}
	}

	// construimos las salidas
	from := string(wallet.GetAddress())
	outputs = append(outputs, *NewTXOuput(amount, to))

	if acc > amount {
		outputs = append(outputs, *NewTXOuput(acc-amount, from))
	}

	tx := Transaction{
		ID:   nil,
		Vin:  inputs,
		Vout: outputs,
	}
	tx.ID = tx.Hash()

	// por ultimo firmamos la transaccion con la clave privada
	UTXOSet.Blockchain.SignTransaction(&tx, wallet.PrivateKey)

	return &tx
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}

	// verificamos si la transacion fue correcta
	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR: Previous transaction not correct")
		}
	}

	// creamos una copia
	txCopy := tx.TrimmedCopy()

	// establecemos los campos de la firma
	for inID, vin := range txCopy.Vin {
		prevTX := prevTXs[hex.EncodeToString(vin.Txid)]

		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTX.Vout[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil

		// establecemos la firma
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)

		if err != nil {
			log.Panic(err)
		}

		tx.Vin[inID].Signature = append(r.Bytes(), s.Bytes()...)
	}
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, TXInput{vin.Txid, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Vout {
		outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash})
	}

	txCopy := Transaction{
		ID:   tx.ID,
		Vin:  inputs,
		Vout: outputs,
	}

	return txCopy
}

// verify the signatures of transaction inputs
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	for _, vin := range tx.Vin {
		if (prevTXs[hex.EncodeToString(vin.Txid)]).ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil

		r, s := big.Int{}, big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x, y := big.Int{}, big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if !ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) {
			return false
		}
	}

	return true
}

// muestra las transacciones realizadas en formato string
func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x", tx.ID))

	for i, input := range tx.Vin {
		lines = append(lines, fmt.Sprintf("	Input %d:", i))
		lines = append(lines, fmt.Sprintf("	TXID:      %x", input.Txid))
		lines = append(lines, fmt.Sprintf("	Out:       %d", input.Vout))
		lines = append(lines, fmt.Sprintf("	Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("	PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.Vout {
		lines = append(lines, fmt.Sprintf("	Output 	%d:", i))
		lines = append(lines, fmt.Sprintf("	Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("	Script: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}
