package p2p

import (
	"dfs/codec"
	"dfs/domain"
	"errors"
	"fmt"
	"net"
	"time"
)

type TCPTransportOpts struct {
	HandshakeFunc HandshakeFunc
	ListenAddress string
	Codec         codec.Codec
	OnPeer        func(peer Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener

	consumeChan chan domain.Message
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		consumeChan:      make(chan domain.Message, 1024),
	}
}

func (t *TCPTransport) Consume() <-chan domain.Message {
	return t.consumeChan
}

func (t *TCPTransport) Close() error {
	fmt.Println("[TCPTransport]: Closing listener")
	return t.listener.Close()
}

func (t *TCPTransport) ListenerAddress() string {
	return t.TCPTransportOpts.ListenAddress
}

func (t *TCPTransport) Dial(addr string) error {
	var conn net.Conn
	var err error

	for {
		conn, err = net.Dial("tcp", addr)
		if err == nil {
			// Connection successful, break the loop
			break
		}

		fmt.Println("[TCPTransport]: Dial error:", err)
		time.Sleep(5 * time.Second)
	}

	go t.handleConn(conn)
	return nil
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddress)
	if err != nil {
		return err
	}

	fmt.Printf("[TCPTransport]: Listening on %s\n", t.ListenAddress)

	go t.startAcceptLoop()
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()

		if errors.Is(err, net.ErrClosed) {
			return
		}

		if err != nil {
			fmt.Println("[TCPTransport]: Connection ACK error:", err)
		}
		fmt.Println("Connection accepted")
		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	fmt.Println("[TCPTransport]: Incoming connection to: ", t.TCPTransportOpts.ListenAddress, "from: ", conn.RemoteAddr().String())

	if err := t.HandshakeFunc(peer); err != nil {
		fmt.Println("[TCPTransport]: Shake hands error:", err)
		return
	}

	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			fmt.Println("[TCPTransport]: OnPeer error:", err)
			return
		}
	}

	msg := &domain.Message{}
	// Read loop
	for {
		if err := t.Codec.Decode(conn, msg); err != nil {
			fmt.Println("[TCPTransport]: Decode error:", err)
			return
		}

		msg.From = conn.RemoteAddr().String()
		msg.Type = domain.MessageTypeStoreFile

		fmt.Println("--------------------------------")
		fmt.Printf("[TCPTransport]: Message Received: [%+v] -> %+v\n", msg.From, msg)
		fmt.Println("--------------------------------")

		peer.Wg.Add(1)
		t.consumeChan <- *msg
		fmt.Println("waiting for peer to finish")
		peer.Wg.Wait()
		fmt.Println("peer finished")
	}
}
