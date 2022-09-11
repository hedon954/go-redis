package database

import (
	"strings"

	"go-redis/datastruct/dict"
	"go-redis/interface/database"
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

// DB represent a redis database
type DB struct {
	index  int
	data   dict.Dict
	addAof func(CmdLine)
}

// makeDB creates the first redis database
func makeDB() *DB {
	db := &DB{
		data:   dict.MakeSyncDict(),
		addAof: func(line CmdLine) {}, // avoid writing aof again while loadAof
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

// GetEntity gets data entity bay key
func (db *DB) GetEntity(key string) (*database.DataEntity, bool) {
	val, exists := db.data.Get(key)
	if !exists {
		return nil, false
	}
	entity, _ := val.(*database.DataEntity)
	return entity, true
}

// PutEntity stores data
func (db *DB) PutEntity(key string, entity *database.DataEntity) int {
	return db.data.Put(key, entity)
}

// PutIfExists stores data if exists
func (db *DB) PutIfExists(key string, entity *database.DataEntity) int {
	return db.data.PutIfExists(key, entity)
}

// PutIfAbsent stores data if not exists
func (db *DB) PutIfAbsent(key string, entity *database.DataEntity) int {
	return db.data.PutIfAbsent(key, entity)
}

// Remove removes a key
func (db *DB) Remove(key string) int {
	return db.data.Remove(key)
}

// Removes removes a list of keys
func (db *DB) Removes(keys ...string) int {
	deleted := 0
	for _, key := range keys {
		deleted += db.data.Remove(key)
	}
	return deleted
}

// Flush clears the DB
func (db *DB) Flush() {
	db.data.Clear()
}

// validateArity checks arity validation
// we have the following appointment:
// SET KV -> arity = 3
// EXISTS k1 k2 k3 k4 ....  arity -> -2 (means more than 2 or equals to 2)
func validateArity(arity int, cmdArgs [][]byte) bool {
	argNum := len(cmdArgs)
	if arity >= 0 {
		return argNum == arity
	}
	return argNum >= -arity
}
