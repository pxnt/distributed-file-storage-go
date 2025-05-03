package p2p

import (
	"fmt"
	"net"
	"sync"
)

type TCPPeer struct {
	conn net.Conn
	// outbound is true if the peer is an outbound peer
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransportOpts struct {
	HandshakeFunc HandshakeFunc
	ListenAddress string
	Decoder       Decoder
}

type TCPTransport struct {
	TCPTransportOpts
	listener      net.Listener

	mu            sync.RWMutex
	peers         map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddress)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		fmt.Println("TCPTransport: Accept Loop")

		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Println("TCPTransport: Accept error:", err)
		}
		fmt.Println("TCPTransport: Accepted Loop")

		go t.handleConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	
	fmt.Println("TCPTransport: New Incoming connection:", peer)

	if err := t.HandshakeFunc(peer); err != nil {
		fmt.Println("TCPTransport: Shake hands error:", err)

	}

	msg := &Message{}
	// Read loop
	for {
		if err := t.Decoder.Decode(conn, msg); err != nil {
			fmt.Println("TCPTransport: Decode error:", err)
			continue
		}

		msg.From = conn.RemoteAddr()

		fmt.Printf("TCPTransport: Decoded message: %+v\n", msg)
	}
}

 