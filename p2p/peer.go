package p2p

import "net"

type Peer interface {
	Send(msg []byte) error
	CloseStream()
	// Close() error
	// RemoteAddress() net.Addr
	net.Conn
}
