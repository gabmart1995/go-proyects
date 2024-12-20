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
