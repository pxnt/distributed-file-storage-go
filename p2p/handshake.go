package p2p

type HandshakeFunc func(Peer) error

func NOPHandshake(peer Peer) error {
	return nil
}
