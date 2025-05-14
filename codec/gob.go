package codec

import (
	"dfs/domain"
	"encoding/gob"
	"io"
)

type GOBCodec struct{}

func (g GOBCodec) Decode(r io.Reader, msg *domain.Message) error {
	return gob.NewDecoder(r).Decode(msg)
}
