package database

import (
	"go-redis/datastruct/dict"
	"go-redis/interface/resp"
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
