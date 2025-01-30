package p2p

import (
	"bittorrent-client/client"
	"bittorrent-client/message"
	"bittorrent-client/peers"
	"bytes"
	"crypto/sha1"
	"fmt"
	"log"
	"runtime"
	"time"
)

// es el largo del numero de bytes en la cual un peticion puede mandar
const MaxBlockSize = 16384

// intentos de conexion
const MaxBackLog = 5

// el torrent mantiene la data requerida desde la lista de peers
type Torrent struct {
	Peers       []peers.Peer
	PeerID      [20]byte
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

type pieceWork struct {
	index  int
	hash   [20]byte
	length int
}

type pieceResult struct {
	index int
	buf   []byte
}

type pieceProgress struct {
	index      int
	client     *client.Client
	buf        []byte
	downloaded int
	requested  int
	backlog    int
}

func (state *pieceProgress) readMessage() error {
	msg, err := state.client.Read() // llamara a los bloques

	if err != nil {
		return err
	}

	if msg == nil {
		return nil
	}

	switch msg.ID {
	case message.MsgUnchoke:
		state.client.Chocked = false

	case message.MsgChoke:
		state.client.Chocked = true

	case message.MsgHave:
		index, err := message.ParseHave(msg)

		if err != nil {
			return err
		}

		state.client.Bitfield.SetPiece(index)

	case message.MsgPiece:
		n, err := message.ParsePiece(state.index, state.buf, msg)

		if err != nil {
			return err
		}

		state.downloaded += n
		state.backlog--
	}

	return nil
}

func (t *Torrent) calculateBoundsForPiece(index int) (begin int, end int) {
	begin = index * t.PieceLength
	end = begin + t.PieceLength

	if end > t.Length {
		end = t.Length
	}

	return begin, end
}

func (t *Torrent) calculatePieceSize(index int) int {
	begin, end := t.calculateBoundsForPiece(index)
	return (end - begin)
}

func checkIntegrity(pw *pieceWork, buf []byte) error {
	hash := sha1.Sum(buf)

	if !bytes.Equal(hash[:], pw.hash[:]) {
		return fmt.Errorf("index %d failed integrity check", pw.index)
	}

	return nil
}

// inicia el worker
func (t *Torrent) startDownloadWorker(peer peers.Peer, workQueue chan *pieceWork, results chan *pieceResult) {
	c, err := client.New(peer, t.PeerID, t.InfoHash)

	if err != nil {
		log.Printf("could not handshake with %s. Disconnecting\n", peer.IP)
		return
	}

	defer c.Conn.Close()

	log.Printf("completed handshake with %s\n", peer.IP)

	c.SendUnchoke()
	c.SendInterested()

	for pw := range workQueue {
		if !c.Bitfield.HasPiece(pw.index) {
			workQueue <- pw // coloca la pieza de vuelta al canal
			continue
		}

		// inicia la descarga de la pieza
		buf, err := attemptDownloadPiece(c, pw)

		if err != nil {
			log.Println("Exiting", err)
			workQueue <- pw // devuelve la pieza
			return
		}

		if err := checkIntegrity(pw, buf); err != nil {
			log.Printf("piece #%d failed integrity check\n", pw.index)
			workQueue <- pw
			continue
		}

		c.SendHave(pw.index)

		results <- &pieceResult{
			index: pw.index,
			buf:   buf,
		}
	}
}

func attemptDownloadPiece(c *client.Client, pw *pieceWork) ([]byte, error) {
	state := pieceProgress{
		index:  pw.index,
		client: c,
		buf:    make([]byte, pw.length),
	}

	// establecemos el limite de conexion
	c.Conn.SetDeadline(time.Now().Add(30 * time.Second))
	defer c.Conn.SetDeadline(time.Time{})

	for state.downloaded < pw.length {
		if !state.client.Chocked {
			for state.backlog < MaxBackLog && state.requested < pw.length {
				blockSize := MaxBlockSize

				// el ultimo bloque debe ser mas corto que los demas
				if (pw.length - state.requested) < blockSize {
					blockSize = pw.length - state.requested
				}

				if err := c.SendRequest(pw.index, state.requested, blockSize); err != nil {
					return nil, err
				}

				state.backlog++
				state.requested += blockSize
			}
		}

		if err := state.readMessage(); err != nil {
			return nil, err
		}
	}

	return state.buf, nil
}

// inicia la descarga del archivo
func (t *Torrent) Download() ([]byte, error) {
	log.Println("Starting download for", t.Name)

	// iniciamos los workers que son canales gestionados por routines
	workQueue := make(chan *pieceWork, len(t.PieceHashes))
	results := make(chan *pieceResult)

	for index, hash := range t.PieceHashes {
		length := t.calculatePieceSize(index)
		workQueue <- &pieceWork{
			index:  index,
			hash:   hash,
			length: length,
		}
	}

	// iniciamos los workers
	for _, peer := range t.Peers {
		go t.startDownloadWorker(peer, workQueue, results)
	}

	buf := make([]byte, t.Length)
	donePieces := 0

	for donePieces < len(t.PieceHashes) {
		res := <-results
		begin, end := t.calculateBoundsForPiece(res.index)

		copy(buf[begin:end], res.buf)
		donePieces++

		// calculamos el porcentaje
		percent := float64(donePieces) / float64(len(t.PieceHashes)) * 100
		numWorkers := (runtime.NumGoroutine() - 1) // obtiene el numero de goroutines menos el principal

		log.Printf(
			"(%0.2f%%) downloaded piece #%d from %d peers\n",
			percent,
			res.index,
			numWorkers,
		)
	}

	close(workQueue)

	return buf, nil
}
