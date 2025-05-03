package main

import (
	"dfs/p2p"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Hello, World!")

	tcpOpts := p2p.TCPTransportOpts{
		ListenAddress: ":3000",
		HandshakeFunc: p2p.NOPHandshake,
		Decoder:       p2p.DefaultDecoder{},
	}
	tr := p2p.NewTCPTransport(tcpOpts)

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
