package main

import (
	"bytes"
	"dfs/domain"
	"dfs/p2p"
	"fmt"
	"io"
	"sync"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOpts

	peerLock sync.Mutex
	peers    map[string]p2p.Peer

	store    *Store
	quitChan chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}

	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitChan:       make(chan struct{}),
		peers:          make(map[string]p2p.Peer),
	}
}

func (s *FileServer) Start() error {
	err := s.Transport.ListenAndAccept()

	if err != nil {
		return err
	}

	s.bootstrapNetwork()

	go s.loop()

	return nil
}

func (s *FileServer) Stop() {
	close(s.quitChan)
}

func (s *FileServer) OnPeer(peer p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[peer.RemoteAddr().String()] = peer

	fmt.Printf("[FileServer]: Peer connected: %s\n", peer.RemoteAddr().String())

	return nil
}

func (s *FileServer) StoreData(key string, r io.Reader) error {
	// 1. store the file to disk
	// 2. broadcast the file to the network

	fmt.Println("\n*****************NEW REQUEST***********************")
	fmt.Printf("[FileServer]: Storing data for key: %s\n", key)

	dataBuffer := new(bytes.Buffer)
	tee := io.TeeReader(r, dataBuffer)

	size, err := s.store.WriteStream(key, tee)

	if err != nil {
		return err
	}

	msg := domain.BroadcastMessage{
		Key:  key,
		Size: size,
	}

	s.broadcastToPeers(&msg)

	// time.Sleep(3 * time.Second)

	s.streamToPeers(dataBuffer.Bytes())

	return nil
}

func (s *FileServer) broadcastToPeers(msg *domain.BroadcastMessage) error {
	buf, err := msg.Encode()
	if err != nil {
		return err
	}

	for _, peer := range s.peers {
		if err := peer.Send(buf); err != nil {
			return err
		}
	}

	fmt.Println("[FileServer]: Message sent to peers")
	return nil
}

func (s *FileServer) streamToPeers(data []byte) (int, error) {
	peers := []io.Writer{}

	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	multiwriter := io.MultiWriter(peers...)

	return multiwriter.Write(data)
}

func (s *FileServer) loop() {
	defer s.Transport.Close()

	for {
		select {
		case msg := <-s.Transport.Consume():
			if err := s.handleMessage(&msg); err != nil {
				fmt.Printf("[FileServer]: Error handling message: %v\n", err)
				return
			}

		case <-s.quitChan:
			fmt.Println("[FileServer]: STOPPED")
			s.Transport.Close()
			return
		}
	}
}

func (s *FileServer) handleMessage(msg *domain.Message) error {
	switch msg.Type {
	case domain.MessageTypeStoreFile:
		return s.handleMessageStoreFile(msg)
	}

	return nil
}

func (s *FileServer) handleMessageStoreFile(msg *domain.Message) error {
	peer, ok := s.peers[msg.From]
	if !ok {
		return fmt.Errorf("[FileServer]: Peer not found: %s\n", msg.From)
	}

	payload, err := msg.DecodePayload()
	if err != nil {
		return err
	}

	fmt.Printf("[FileServer]: ingesting stream of size: %+v\n", payload.Size)

	s.store.WriteStream(payload.Key, io.LimitReader(peer, payload.Size))

	peer.CloseStream()

	fmt.Printf("[FileServer]: stream ingested")
	return nil
}

func (s *FileServer) bootstrapNetwork() {
	wg := sync.WaitGroup{}

	for _, addr := range s.BootstrapNodes {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			fmt.Printf("[FileServer]: Dialing bootstrap node: %s\n", addr)
			err := s.Transport.Dial(addr)
			if err != nil {
				fmt.Printf("[FileServer]: Bootstrap node error: %v\n", err)
			}
		}(addr)
	}
	wg.Wait()
}

func (s *FileServer) Get(key string) (io.Reader, error) {
	if s.store.Has(key) {
		return s.store.ReadStream(key)
	}

	return nil, fmt.Errorf("[FileServer]: Key not found: %s\n", key)
}
