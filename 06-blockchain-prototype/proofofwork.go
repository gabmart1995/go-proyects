package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var maxNonce = math.MaxInt64

const targetBits = 24 // difulcultad de la mineria

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	// gereneramos un big int con el valor por defecto de 1
	target := big.NewInt(1)

	// desplazamos el valor hacia la izquierda 256 bits quedando la direccion
	// 0x10000000000000000000000000000000000... usado en el algoritmo SHA256
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{
		block:  b,
		target: target,
	}

	return pow
}

// prepara la prueba de trabajo
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte

	nonce := 0 // intentos para conseguir el hash correcto

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)

	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)

		// establece los bytes en el big
		hashInt.SetBytes(hash[:])

		// comprobamos si cumple con la prueba de trabajo
		if hashInt.Cmp(pow.target) == -1 {
			break

		} else {
			nonce++

		}
	}

	fmt.Printf("\n\n")

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)

	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
