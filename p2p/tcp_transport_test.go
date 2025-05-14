package p2p

import (
	"dfs/codec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTCPTransport(t *testing.T) {
	tcpOpts := TCPTransportOpts{
		ListenAddress: ":3000",
		HandshakeFunc: NOPHandshake,
		Codec:         codec.DefaultCodec{},
	}
	tr := NewTCPTransport(tcpOpts)

	assert.Equal(t, tr.ListenAddress, ":3000")

	assert.Nil(t, tr.ListenAndAccept())
}
