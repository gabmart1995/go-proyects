package client

import (
	"bittorrent-client/bitfield"
	"bittorrent-client/peers"
	"net"
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
/*func New(peer peers.Peer, peerId, infoHash [20]byte) (*Client, error) {
	conn, err := net.DialTimeout("tcp", peer.String(), 3*time.Second)

	if err != nil {
		return nil, err
	}

}*/

func completeHandshake(conn net.Conn, infoHash, peerID [20]byte) {

}
