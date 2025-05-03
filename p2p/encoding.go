package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *Message) error
}

type DefaultDecoder struct{}

func (d DefaultDecoder) Decode(r io.Reader, msg *Message) error {
	buf := make([]byte, 1024)

	n, err := r.Read(buf)
	if err != nil {
		return err
	}

	msg.Payload = buf[:n]
	return nil
}

type GOBDecoder struct{}

func (g GOBDecoder) Decode(r io.Reader, msg *Message) error {
	return gob.NewDecoder(r).Decode(msg)
}