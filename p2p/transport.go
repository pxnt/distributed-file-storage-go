package p2p

import (
	"dfs/domain"
)

type Transport interface {
	ListenerAddress() string
	Dial(addr string) error
	ListenAndAccept() error
	Consume() <-chan domain.Message
	Close() error
}
