package domain

import (
	"bytes"
	"encoding/gob"
)

// Message to be sent to peers
type BroadcastMessage struct {
	Key  string // file name
	Size int64  // data bytes size
}

func (b BroadcastMessage) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := gob.NewEncoder(buf).Encode(b); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type MessageType int

const (
	MessageTypeStoreFile MessageType = iota
	MessageTypeStream
)

// Message received to peers
type Message struct {
	Type    MessageType
	From    string
	Payload []byte // BroadcastMessage
}

func (m Message) DecodePayload() (*BroadcastMessage, error) {
	var payload BroadcastMessage
	if err := gob.NewDecoder(bytes.NewReader(m.Payload)).Decode(&payload); err != nil {
		return nil, err
	}

	return &payload, nil
}
