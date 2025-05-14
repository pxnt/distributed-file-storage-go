package codec

import (
	"dfs/domain"
	"io"
)

type DefaultCodec struct{}

func (d DefaultCodec) Decode(r io.Reader, msg *domain.Message) error {
	buf := make([]byte, 1024)

	n, err := r.Read(buf)
	if err != nil {
		return err
	}

	msg.Payload = buf[:n]
	return nil
}
