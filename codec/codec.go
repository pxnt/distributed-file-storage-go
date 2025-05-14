package codec

import (
	"dfs/domain"
	"io"
)

type Codec interface {
	Decode(io.Reader, *domain.Message) error
}
