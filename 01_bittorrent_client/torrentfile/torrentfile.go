package torrentfile

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"

	"github.com/jackpal/bencode-go"
)

// puerto de conexion a los bittorrent
const PORT uint16 = 6881

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

// lee el contenido del archivo torrent
func Open(r io.Reader) (*bencodeTorrent, error) {
	bto := bencodeTorrent{}

	if err := bencode.Unmarshal(r, &bto); err != nil {
		return nil, err
	}

	return &bto, nil
}

// genera el hash para el bencode
func (i *bencodeInfo) hash() ([20]byte, error) {
	var buf bytes.Buffer

	if err := bencode.Marshal(&buf, *i); err != nil {
		return [20]byte{}, err
	}

	h := sha1.Sum(buf.Bytes())

	return h, nil
}

// retorna un array de los hashes peers
func (i *bencodeInfo) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20 // longitud del sha-1 hash
	buf := []byte(i.Pieces)

	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("received malformed pieces of length %d", len(buf))
		return nil, err
	}

	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)

	// copia una seccion del buffer y lo asigna en el array
	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[(i*hashLen):((i+1)*hashLen)])
	}

	return hashes, nil
}

func (bto *bencodeTorrent) toTorrentFile() (TorrentFile, error) {
	infoHash, err := bto.Info.hash()

	if err != nil {
		return TorrentFile{}, err
	}

	pieceHashes, err := bto.Info.splitPieceHashes()

	if err != nil {
		return TorrentFile{}, err
	}

	t := TorrentFile{
		Announce:    bto.Announce,
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: bto.Info.PieceLength,
		Length:      bto.Info.Length,
		Name:        bto.Info.Name,
	}

	return t, nil
}
