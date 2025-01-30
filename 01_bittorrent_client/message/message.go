package message

import (
	"encoding/binary"
	"fmt"
	"io"
)

type messageID uint8

// estados del mensaje protocolo bittorrent
const (
	// MsgChoke chokes the receiver
	MsgChoke messageID = 0
	// MsgUnchoke unchokes the receiver
	MsgUnchoke messageID = 1
	// MsgInterested expresses interest in receiving data
	MsgInterested messageID = 2
	// MsgNotInterested expresses disinterest in receiving data
	MsgNotInterested messageID = 3
	// MsgHave alerts the receiver that the sender has downloaded a piece
	MsgHave messageID = 4
	// MsgBitfield encodes which pieces that the sender has downloaded
	MsgBitfield messageID = 5
	// MsgRequest requests a block of data from the receiver
	MsgRequest messageID = 6
	// MsgPiece delivers a block of data to fulfill a request
	MsgPiece messageID = 7
	// MsgCancel cancels a request
	MsgCancel messageID = 8
)

type Message struct {
	ID      messageID
	Payload []byte
}

// crea un mensaje de peticion de recursos
func FormatRequest(index, begin, length int) *Message {
	payload := make([]byte, 12)

	// los payload se manejan en bigendian 32bits
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))

	return &Message{
		ID:      MsgRequest,
		Payload: payload,
	}
}

// crea un mesaje verificacion del payload
func FormatHave(index int) *Message {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, uint32(index))

	return &Message{
		ID:      MsgHave,
		Payload: payload,
	}
}

// extrae los datos del buffer del mensaje Piece
func ParsePiece(index int, buf []byte, msg *Message) (int, error) {
	if msg.ID != MsgPiece {
		return 0, fmt.Errorf("expected piece (ID %d), got ID %d", MsgPiece, msg.ID)
	}

	// verificamos la longitud
	if len(msg.Payload) < 8 {
		return 0, fmt.Errorf("payload too short. %d < 8", len(msg.Payload))
	}

	parsedIndex := int(binary.BigEndian.Uint32(msg.Payload[0:4]))

	if parsedIndex != index {
		return 0, fmt.Errorf("expected index %d, got %d", index, parsedIndex)
	}

	begin := int(binary.BigEndian.Uint32(msg.Payload[4:8]))

	if begin >= len(buf) {
		return 0, fmt.Errorf("begin offset too high. %d >= %d", begin, len(buf))
	}

	data := msg.Payload[8:]

	if (begin + len(data)) > len(buf) {
		return 0, fmt.Errorf(
			"data too long [%d] for offset %d with length %d",
			len(data),
			begin,
			len(buf),
		)
	}

	// copiamos los datos al archivo
	copy(buf[begin:], data)

	return len(data), nil
}

// extrae los datos del mensaje Have
func ParseHave(msg *Message) (int, error) {
	if msg.ID != MsgHave {
		return 0, fmt.Errorf("expected HAVE (ID %d), got ID %d", MsgHave, msg.ID)
	}

	if len(msg.Payload) != 4 {
		return 0, fmt.Errorf("expected payload length 4, got length %d", len(msg.Payload))
	}

	index := int(binary.BigEndian.Uint32(msg.Payload))

	return index, nil
}

// serializa el mensaje hacia el buffer
// <length prefix><message ID><payload>
func (m *Message) Serialize() []byte {
	if m == nil {
		return make([]byte, 4)
	}

	length := uint32(len(m.Payload) + 1) // +1 para ID
	buf := make([]byte, 4+length)

	binary.BigEndian.PutUint32(buf[0:4], length)

	// a partir de la posicion 4 se obtiene el ID
	buf[4] = byte(m.ID)

	copy(buf[5:], m.Payload)

	return buf
}

// parse el mensaje desde el stream.
func Read(r io.Reader) (*Message, error) {
	lengthBuf := make([]byte, 4)

	if _, err := io.ReadFull(r, lengthBuf); err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(lengthBuf)

	// en caso de llegar vacio mantiene la conexion
	if length == 0 {
		return nil, nil
	}

	messageBuf := make([]byte, length)

	if _, err := io.ReadFull(r, messageBuf); err != nil {
		return nil, err
	}

	m := Message{
		ID:      messageID(messageBuf[0]),
		Payload: messageBuf[1:],
	}

	return &m, nil
}

// obtiene el estado de la conversación peer to peer
func (m *Message) name() string {
	if m == nil {
		return "KeepAlive"
	}

	switch m.ID {
	case MsgChoke:
		return "Choke"
	case MsgUnchoke:
		return "Unchoke"
	case MsgInterested:
		return "Interested"
	case MsgNotInterested:
		return "NotInterested"
	case MsgHave:
		return "Have"
	case MsgBitfield:
		return "Bitfield"
	case MsgRequest:
		return "Request"
	case MsgPiece:
		return "Piece"
	case MsgCancel:
		return "Cancel"
	default:
		return fmt.Sprintf("Unknown#%d", m.ID)
	}
}

func (m *Message) String() string {
	if m == nil {
		return m.name()
	}

	return fmt.Sprintf("%s [%d]", m.name(), len(m.Payload))
}
