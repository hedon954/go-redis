package tcp

import (
	"bufio"
	"context"
	"io"
	"net"
	"sync"
	"time"

	"go-redis/lib/logger"
	"go-redis/lib/sync/atomic"
	"go-redis/lib/sync/wait"
)

// EchoHandler implements tcp.Handler interface
type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean
}

type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

func MakeHandler() *EchoHandler {
	return &EchoHandler{}
}

// Close implements `closer` interface
func (client *EchoClient) Close() error {
	client.Waiting.WaitWithTimeout(10 * time.Second)
	_ = client.Conn.Close() // ignore error
	return nil
}

// Handle handles the client connections
func (handler *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	// server is closing, refuse conn
	if handler.closing.Get() {
		_ = conn.Close()
	}

	// accept conn
	client := &EchoClient{
		Conn: conn,
	}

	// value: struct{}{}, use `zerobase` to save space
	handler.activeConn.Store(client, struct{}{})

	// serve conn
	reader := bufio.NewReader(conn)
	for true {
		// receive msg
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF { // client close, msg end
				logger.Info("connection closed")
				handler.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}

		// send msg back(optional)
		client.Waiting.Add(1)
		bs := []byte(msg)
		_, _ = conn.Write(bs)
		client.Waiting.Done()
	}
}

// Close closes handler
func (handler *EchoHandler) Close() error {
	logger.Info("handler shutting down")

	// set closing flag
	handler.closing.Set(true)

	// close all connections
	handler.activeConn.Range(func(key, value interface{}) bool {
		ec := key.(*EchoClient)
		_ = ec.Conn.Close()
		return true // return true means to handle next conn
	})

	return nil
}
