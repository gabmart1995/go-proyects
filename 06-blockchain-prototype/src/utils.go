package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

func IntToHex(num int64) []byte {
	var buffer bytes.Buffer

	if err := binary.Write(&buffer, binary.BigEndian, num); err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}

// invierte las posiciones un array de bytes
func ReverseBytes(data []byte) {
	// bucle de 2 variables iteradoras
	for i, j := 0, (len(data) - 1); i < j; i, j = (i + 1), (j - 1) {
		data[i], data[j] = data[j], data[i]
	}
}
