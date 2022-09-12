package cluster

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

// renameFunc RENAME k1 k2
func renameFunc(cdb *Database, c resp.Connection, cmdArgs [][]byte) resp.Reply {

	if len(cmdArgs) != 3 {
		return reply.MakeArgNumErrReply("rename")
	}

	src := string(cmdArgs[1])
	dst := string(cmdArgs[2])

	srcNode := cdb.peerPicker.PickNode(src)
	dstNode := cdb.peerPicker.PickNode(dst)

	if srcNode != dstNode {
		return reply.MakeStandardErrReply("ERR rename must within on redis node")
	}

	return cdb.relay(srcNode, c, cmdArgs)
}
