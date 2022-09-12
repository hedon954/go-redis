package cluster

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

func delFunc(cdb *Database, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	replies := cdb.broadcast(c, cmdArgs)

	var delCount int64 = 0

	var errReply reply.ErrorReply
	for _, r := range replies {
		if reply.IsErrReply(r) {
			errReply = r.(reply.ErrorReply)
			break
		}
		intReply, ok := r.(*reply.IntReply)
		if !ok {
			errReply = reply.MakeStandardErrReply("error")
			break
		}
		delCount += intReply.Code
	}

	if errReply != nil {
		return reply.MakeStandardErrReply("error: " + errReply.Error())
	}
	
	return reply.MakeIntReply(delCount)
}
