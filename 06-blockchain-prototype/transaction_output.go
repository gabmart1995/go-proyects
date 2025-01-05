package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

type TXOutput struct {
	Value      int
	PubKeyHash []byte
}

type TXOutputs struct {
	Outputs []TXOutput
}

func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]

	out.PubKeyHash = pubKeyHash
}

func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Equal(out.PubKeyHash, pubKeyHash)
}

func NewTXOuput(value int, address string) *TXOutput {
	txo := &TXOutput{
		Value:      value,
		PubKeyHash: nil,
	}

	txo.Lock([]byte(address))

	return txo
}

// serializa las salida de las transacciones
func (outs TXOutputs) Serialize() []byte {
	var buff bytes.Buffer

	encoder := gob.NewEncoder(&buff)

	if err := encoder.Encode(outs); err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func DeserializeOutputs(data []byte) TXOutputs {
	var outputs TXOutputs

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&outputs)

	if err != nil {
		log.Panic(err)
	}

	return outputs
}
