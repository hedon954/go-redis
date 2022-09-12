package handler

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	"go-redis/database"
	"go-redis/lib/logger"
	"go-redis/lib/sync/atomic"
	"go-redis/resp/connection"
	"go-redis/resp/parser"
	"go-redis/resp/reply"

	databaseface "go-redis/interface/database"
)

// RespHandler handlers information that complies with the RESP protocol
type RespHandler struct {
	activeConn sync.Map
	db         databaseface.Database
	closing    atomic.Boolean
}

func MakeRespHandler() *RespHandler {
	return &RespHandler{
		db: database.NewStandaloneDatabase(),
	}
}

// closeClient closes specified client connection
func (r *RespHandler) closeClient(client *connection.Connection) {
	_ = client.Close()
	_ = r.db.AfterClientClose(client)
	r.activeConn.Delete(client)
}

// Handle handles client request
func (r *RespHandler) Handle(ctx context.Context, conn net.Conn) {
	// check handler status
	// handler is closing, refuses connection
	if r.closing.Get() {
		_ = conn.Close()
	}

	// adds connection to map
	client := connection.NewConn(conn)
	r.activeConn.Store(client, struct{}{})

	// server connection
	ch := parser.ParseStream(conn)

	// listen to the channel to get handle result
	for payload := range ch {

		fmt.Printf("got: %v\n", payload)

		// error
		if payload.Err != nil {

			if payload.Err == io.EOF ||
				payload.Err == io.ErrUnexpectedEOF ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") {
				r.closeClient(client)
				logger.Info(fmt.Sprintf("connection closed: %v", client.RemoteAddr()))
				return
			}

			errReply := reply.MakeStandardErrReply(payload.Err.Error())
			err := client.Write(errReply.ToBytes())
			if err != nil {
				r.closeClient(client)
				logger.Info("connection closed: %v", client.RemoteAddr())
				return
			}

			continue
		}

		// exec
		if payload.Data == nil {
			logger.Error("empty payload")
			continue
		}

		bulkReply, ok := payload.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("require multiBulkReplay")
			continue
		}

		result := r.db.Exec(client, bulkReply.Args)
		if result != nil {
			_ = client.Write(result.ToBytes())
		} else {
			_ = client.Write(reply.MakeUnknownErrReplay().ToBytes())
		}
	}
}

// Close closes handler, database and the connections
func (r *RespHandler) Close() error {
	logger.Info("handler shutting down...")
	r.closing.Set(true)
	r.activeConn.Range(func(key, value interface{}) bool {
		client := key.(*connection.Connection)
		_ = client.Close()
		return true
	})
	_ = r.db.Close()
	return nil
}
