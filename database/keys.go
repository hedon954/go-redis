package database

import (
	"go-redis/interface/resp"
	"go-redis/lib/utils"
	"go-redis/lib/wildcard"
	"go-redis/resp/reply"
)

func init() {
	RegisterCommand("del", execDel, -2)
	RegisterCommand("exists", execExists, -2)
	RegisterCommand("flushdb", execFlushDB, -1)
	RegisterCommand("type", execType, 2)
	RegisterCommand("rename", execRename, 3)
	RegisterCommand("renamenx", execRenameNX, 3)
	RegisterCommand("keys", execKeys, 2)
}

// execDel DEL k1 k2 k3 ...
func execDel(db *DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}
	deleted := db.Removes(keys...)
	if deleted > 0 {
		db.addAof(utils.ToCmdLine2("del", args...))
	}
	return reply.MakeIntReply(int64(deleted))
}

// execExists EXISTS k1 k2 k3 k4 ...
func execExists(db *DB, args [][]byte) resp.Reply {
	result := int64(0)
	for _, arg := range args {
		key := string(arg)
		_, exists := db.GetEntity(key)
		if exists {
			result++
		}
	}
	return reply.MakeIntReply(result)
}

// execFlushDB FLUSHDB
func execFlushDB(db *DB, args [][]byte) resp.Reply {
	db.data.Clear()
	db.addAof(utils.ToCmdLine2("flushdb", args...))
	return reply.MakeOKReply()
}

// execType TYPE k1
func execType(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.MakeStatusReply("none") // :none\r\n
	}
	switch entity.Data.(type) {
	case []byte:
		return reply.MakeStatusReply("string")
	}
	return reply.MakeUnknownErrReplay()
}

// execRename RENAME k1 k2
func execRename(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dst := string(args[1])
	val, exist := db.GetEntity(src)
	if !exist {
		return reply.MakeStatusReply("no such key")
	}
	db.Remove(src)
	db.PutEntity(dst, val)
	db.addAof(utils.ToCmdLine2("rename", args...))
	return reply.MakeOKReply()
}

// execRenameNX RENAMENX k1 k2
func execRenameNX(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dst := string(args[1])
	_, exist := db.GetEntity(dst)
	if exist {
		return reply.MakeIntReply(0)
	}
	val, exist := db.GetEntity(src)
	if !exist {
		return reply.MakeStatusReply("no such key")
	}
	db.Remove(src)
	db.PutEntity(dst, val)
	db.addAof(utils.ToCmdLine2("renamenx", args...))
	return reply.MakeIntReply(1)
}

// execKeys KEYS *
func execKeys(db *DB, args [][]byte) resp.Reply {
	pattern := wildcard.CompilePattern(string(args[0]))
	result := make([][]byte, 0)
	db.data.Foreach(func(key string, val interface{}) bool {
		if pattern.IsMatch(key) {
			result = append(result, []byte(key))
		}
		return true
	})
	return reply.MakeMultiBulkReply(result)
}
