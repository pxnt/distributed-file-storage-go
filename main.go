package main

import (
	"bytes"
	"dfs/codec"
	"dfs/p2p"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

func createServer(listenAddr string, bootstrapNodes ...string) *FileServer {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddress: listenAddr,
		HandshakeFunc: p2p.NOPHandshake,
		Codec:         codec.DefaultCodec{},
	}
	tr := p2p.NewTCPTransport(tcpOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       fmt.Sprintf("%s_network", listenAddr),
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tr,
		BootstrapNodes:    bootstrapNodes,
	}

	fileServer := NewFileServer(fileServerOpts)

	tr.OnPeer = fileServer.OnPeer

	return fileServer
}

// ============= TEST SERVERS  - START =============
var a = 10

func createServerWithBootstrap(key string, val string, listenAddr string, bootstrapNodes ...string) {
	server := createServer(listenAddr, bootstrapNodes...)

	server.Start()

	fmt.Println("\n==================INTERCOM=========================\n")

	time.Sleep(time.Duration(a) * time.Second)

	data := bytes.NewBuffer([]byte(val))
	server.StoreData(key, data)
}

func createTestServers() {
	serverNo := os.Args[1]

	fmt.Printf("Server %s starting...\n", serverNo)

	if serverNo == "s1" {
		createServerWithBootstrap("test1", "meow meow meow", ":3000", ":4000")
	} else if serverNo == "s2" {
		createServerWithBootstrap("test2", "bhow bhow bhow", ":4000")
	}
}

// ============= TEST SERVERS - END =============

// ============= MAIN - START =============

func createMultiServers() {
	s1 := createServer(":3000", ":4000")
	s2 := createServer(":4000")

	wg := sync.WaitGroup{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		s1.Start()
	}()
	go func() {
		defer wg.Done()
		s2.Start()
	}()

	wg.Wait()

	// for better logging
	time.Sleep(3 * time.Second)

	fmt.Println("\n==================INTERCOM=========================\n")

	data := bytes.NewBuffer([]byte("meow meow meow"))
	s1.StoreData("test_path", data)

	fmt.Println("*********CHECKING********")
	time.Sleep(7 * time.Second)

	r, err := s2.Get("test_path")
	if err != nil {
		panic(err)
	}

	a, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(a))
}

// ============= MAIN - END =============

func main() {
	fmt.Println("Initializing DFS...")
	fmt.Println("--------------------------------")

	var serverNo string

	if len(os.Args) > 1 {
		serverNo = os.Args[1]
	}

	if serverNo == "s1" || serverNo == "s2" {
		createTestServers()
	} else {
		createMultiServers()
	}

	select {}
}
