package database

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

func init() {
	RegisterCommand("ping", ping, 1)
}

// ping PING
func ping(db *DB, args [][]byte) resp.Reply {
	return reply.MakePongReply()
}
