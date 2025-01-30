package client

import (
	"bittorrent-client/bitfield"
	"bittorrent-client/handshake"
	"bittorrent-client/peers"
	"bytes"
	"fmt"
	"net"
	"time"
)

type Client struct {
	Conn     net.Conn
	Chocked  bool
	Bitfield bitfield.Bitfield
	peer     peers.Peer
	infoHash [20]byte
	peerID   [20]byte
}

// realiza la conexion con el compañero realiza el handshake y recibe la respuesta
/*func New(peer peers.Peer, peerID, infoHash [20]byte) (*Client, error) {
	conn, err := net.DialTimeout("tcp", peer.String(), 3*time.Second)

	if err != nil {
		return nil, err
	}

	if _, err := completeHandshake(conn, infoHash, peerID); err != nil {
		conn.Close() // cierra la conexion en caso de falla
		return nil, err
	}

}*/

func completeHandshake(conn net.Conn, infoHash, peerID [20]byte) (*handshake.Handshake, error) {
	// establece un tiempo limite para la conexion
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer conn.SetDeadline(time.Time{})

	req := handshake.New(infoHash, peerID)

	// escribmos los datos del mensaje
	if _, err := conn.Write(req.Serialize()); err != nil {
		return nil, err
	}

	// esperamos la respuesta del peer en la conexion
	res, err := handshake.Read(conn)

	if err != nil {
		return nil, err
	}

	// verificamos la integradad del mensaje
	if !bytes.Equal(res.InfoHash[:], infoHash[:]) {
		err = fmt.Errorf("expected infohash %x but got %x", res.InfoHash, infoHash)
		return nil, err
	}

	return res, nil
}

// manejador de los mensajes bitfield
/*func recvBitField(conn net.Conn) (bitfield.Bitfield, error) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{})
}*/
