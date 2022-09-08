package database

import (
	"go-redis/interface/resp"
	"go-redis/lib/logger"
	"go-redis/resp/reply"
)

type EchoDatabase struct {
}

func MakeEchoDatabase() *EchoDatabase {
	return &EchoDatabase{}
}

func (e *EchoDatabase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	bulkReply := reply.MakeMultiBulkReply(args)
	return bulkReply
}

func (e *EchoDatabase) Close() error {
	logger.Info("EchoDatabase AfterClientClose")
	return nil
}

func (e *EchoDatabase) AfterClientClose(client resp.Connection) error {
	logger.Info("EchoDatabase Close")
	return nil
}
