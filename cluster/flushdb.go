package cluster

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

func flushdbFunc(cdb *Database, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	replies := cdb.broadcast(c, cmdArgs)

	var errReply reply.ErrorReply

	for _, r := range replies {
		if reply.IsErrReply(r) {
			errReply = r.(reply.ErrorReply)
			return reply.MakeStandardErrReply("error: " + errReply.Error())
		}
	}

	return reply.MakeOKReply()
}
