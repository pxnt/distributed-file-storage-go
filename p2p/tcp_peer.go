package p2p

import (
	"net"
	"sync"
)

type TCPPeer struct {
	net.Conn
	// outbound is true if the peer is an outbound peer
	outbound bool

	Wg *sync.WaitGroup
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
		Wg:       &sync.WaitGroup{},
	}
}

func (peer *TCPPeer) Send(msg []byte) error {
	_, err := peer.Conn.Write(msg)
	return err
}

func (peer *TCPPeer) CloseStream() {
	peer.Wg.Done()
}
