package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"log"

	"golang.org/x/crypto/ripemd160"
)

/*
	Una wallet es sencillamente una aplicacion que contiene
	las claves publicas y privadas al momento de crear transacciones
*/

const version = byte(0x00)
const addressChecksumLen = 4
const walletFile = "wallet.json"

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// genera un nuevo wallet
func NewWallet() *Wallet {
	private, public := newKeyPair()

	wallet := Wallet{
		PrivateKey: private,
		PublicKey:  public,
	}

	return &wallet
}

// genera el algoritmo de curva elliptica para generar
// las claves del wallet
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)

	if err != nil {
		log.Panic(err)
	}

	pubKey := append(
		private.X.Bytes(),
		private.Y.Bytes()...,
	)

	return *private, pubKey
}

// obtiene la clave publica del wallet
func (w Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)
	versionPayload := append([]byte{version}, pubKeyHash...)

	checksum := checksum(versionPayload)
	fullPayload := append(versionPayload, checksum...)

	address := Base58Encode(fullPayload)

	return address
}

// cifra la clave publica
func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()

	_, err := RIPEMD160Hasher.Write(publicSHA256[:])

	if err != nil {
		log.Panic(err)
	}

	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

// genera un checksum para una clave publica
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

// transforma la informacion del wallet en un JSON
func (w Wallet) MarshalJSON() ([]byte, error) {
	mapStringAny := map[string]any{
		"PrivateKey": map[string]any{
			"D": w.PrivateKey.D,
			"PublicKey": map[string]any{
				"X": w.PrivateKey.PublicKey.X,
				"Y": w.PrivateKey.PublicKey.Y,
			},
			"X": w.PrivateKey.X,
			"Y": w.PrivateKey.Y,
		},
		"PublicKey": w.PublicKey,
	}
	return json.Marshal(mapStringAny)
}
