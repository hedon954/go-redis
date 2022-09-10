package database

import (
	"go-redis/datastruct/dict"
	"go-redis/interface/resp"
	"go-redis/resp/reply"
	"strings"
)

// DB represent a redis database
type DB struct {
	index int
	data  dict.Dict
}

// makeDB creates the first redis database
func makeDB() *DB {
	db := &DB{
		data: dict.MakeSyncDict(),
	}
	return db
}

// ExecFun used to warp redis command
type ExecFun func(db *DB, args [][]byte) resp.Reply

// CmdLine redis command
type CmdLine [][]byte

// Exec runs redis command
func (db *DB) Exec(c resp.Connection, cmdLine CmdLine) resp.Reply {
	if len(cmdTable) <= 0 {
		return reply.MakeStandardErrReply("ERR empty cmd")
	}
	cmdName := strings.ToLower(string(cmdLine[0])) // PING, SET, GET, etc..
	cmd, ok := cmdTable[cmdName]
	if !ok {
		return reply.MakeStandardErrReply("ERR unknown command " + cmdName)
	}

	if !validateArity(cmd.arity, cmdLine) {
		return reply.MakeArgNumErrReply(cmdName) // SET key
	}
	return cmd.executor(db, cmdLine[1:]) // SET k v -> k v
}

// validateArity checks arity validation
func validateArity(arity int, cmdArgs [][]byte) bool {
	return true
}
