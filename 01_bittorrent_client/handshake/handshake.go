package handshake

import (
	"fmt"
	"io"
)

// es un mensaje especial que utilza el peer para identificarse.
type Handshake struct {
	Pstr     string
	InfoHash [20]byte
	PeerID   [20]byte
}

func New(infoHash, peerID [20]byte) *Handshake {
	return &Handshake{
		Pstr:     "BitTorrent protocol",
		InfoHash: infoHash,
		PeerID:   peerID,
	}
}

// serializa el mensaje hacia el buffer
func (h *Handshake) Serialize() []byte {
	buf := make([]byte, len(h.Pstr)+49)
	buf[0] = byte(len(h.Pstr))

	// se realiza una copia incremental
	curr := 1
	curr += copy(buf[curr:], []byte(h.Pstr))
	curr += copy(buf[curr:], make([]byte, 8)) // reserva 8 bits en 0
	curr += copy(buf[curr:], h.InfoHash[:])
	curr += copy(buf[curr:], h.PeerID[:])

	return buf
}

// parsea el handshake desde el stream
func Read(r io.Reader) (*Handshake, error) {
	lengthBuf := make([]byte, 1)

	if _, err := io.ReadFull(r, lengthBuf); err != nil {
		return nil, err
	}

	// verificamos el protocolo
	pstrLen := int(lengthBuf[0])

	if pstrLen == 0 {
		err := fmt.Errorf("pstrlen cannot be 0")
		return nil, err
	}

	// generamos un buffer para almacenar el handshake
	handshakeBuf := make([]byte, 48+pstrLen)

	if _, err := io.ReadFull(r, handshakeBuf); err != nil {
		return nil, err
	}

	var infoHash, peerID [20]byte

	// copiamos los datos del buffer hacia las variables
	copy(infoHash[:], handshakeBuf[pstrLen+8:pstrLen+8+20])
	copy(peerID[:], handshakeBuf[pstrLen+8+20:])

	h := Handshake{
		Pstr:     string(handshakeBuf[0:pstrLen]),
		InfoHash: infoHash,
		PeerID:   peerID,
	}

	return &h, nil
}
