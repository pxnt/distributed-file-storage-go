package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTCPTransport(t *testing.T) {
	tcpOpts := TCPTransportOpts{
		ListenAddress: ":3000",
		HandshakeFunc: NOPHandshake,
		Decoder:       DefaultDecoder{},
	}
	tr := NewTCPTransport(tcpOpts)

	assert.Equal(t, tr.ListenAddress, ":3000")

	assert.Nil(t, tr.ListenAndAccept())
}
