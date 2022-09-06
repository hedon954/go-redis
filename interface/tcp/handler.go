package tcp

import (
	"context"
	"net"
)

// Handler defines the tcp server handler interface
type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}
