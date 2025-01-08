package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

/*
	Una wallet es sencillamente una aplicacion que contiene
	las claves publicas y privadas al momento de crear transacciones
*/

const version = byte(0x00)
const addressChecksumLen = 4

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

type _PrivateKey struct {
	D          *big.Int
	PublicKeyX *big.Int
	PublicKeyY *big.Int
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
		private.PublicKey.X.Bytes(),
		private.PublicKey.Y.Bytes()...,
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

// valida si la direccion es valida
func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Equal(actualChecksum, targetChecksum)
}

// funcion que permite crear una interfaz para la clave privada
// para que luego pueda ser endodeada por el GOB
func (w *Wallet) GobEncode() ([]byte, error) {
	privateKey := &_PrivateKey{
		D:          w.PrivateKey.D,
		PublicKeyX: w.PrivateKey.PublicKey.X,
		PublicKeyY: w.PrivateKey.PublicKey.Y,
	}

	var buf bytes.Buffer
	var err error

	encoder := gob.NewEncoder(&buf)

	if err = encoder.Encode(&privateKey); err != nil {
		return nil, err
	}

	_, err = buf.Write(w.PublicKey)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// metodo que establece las claves dentro del archivo del wallet
// siendo decodeada por el GOB
func (w *Wallet) GobDecode(data []byte) error {
	var privKey _PrivateKey

	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&privKey)

	if err != nil {
		return nil
	}

	// establecemos los campos en el wallet
	w.PrivateKey = ecdsa.PrivateKey{
		D: privKey.D,
		PublicKey: ecdsa.PublicKey{
			X:     privKey.PublicKeyX,
			Y:     privKey.PublicKeyY,
			Curve: elliptic.P256(), // generamos una curva
		},
	}

	// creamos el campo byte con la longitud del buffer
	w.PublicKey = make([]byte, buf.Len())
	_, err = buf.Read(w.PublicKey)

	if err != nil {
		return err
	}

	return nil
}
