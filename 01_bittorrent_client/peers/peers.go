package peers

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

type Peer struct {
	IP   net.IP
	Port uint16
}

// transforma los datos peers de binario en una IP y puerto validos
func Unmarshal(peersBin []byte) ([]Peer, error) {
	const peerSize = 6 // 4 para la IP, 2 para el puerto

	numPeers := len(peersBin) / peerSize

	if len(peersBin)%peerSize != 0 {
		err := fmt.Errorf("receive malformed peers")
		return nil, err
	}

	peers := make([]Peer, numPeers)

	// procedemos a separar la cadena 4 para IP y 2 para el puerto
	for i := 0; i < numPeers; i++ {
		offset := i * peerSize

		peers[i].IP = net.IP(peersBin[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16([]byte(peersBin[offset+4 : offset+6]))
	}

	return peers, nil
}

// transforma los datos en una url valida
func (p Peer) String() string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}
