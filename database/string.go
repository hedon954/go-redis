package database

import (
	"go-redis/interface/database"
	"go-redis/interface/resp"
	"go-redis/lib/utils"
	"go-redis/resp/reply"
)

func init() {
	RegisterCommand("get", execGet, 2)
	RegisterCommand("set", execSet, 3)
	RegisterCommand("setnx", execSetNX, 3)
	RegisterCommand("getset", execGetSet, 3)
	RegisterCommand("strlen", execStrLen, 2)
}

// execGet GET k1
func execGet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exist := db.GetEntity(key)
	if !exist {
		return reply.MakeNullBulkReply()
	}
	return reply.MakeBulkReply(entity.Data.([]byte))
}

// execSet SET k v
func execSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := args[1]
	entity := &database.DataEntity{
		Data: val,
	}

	db.PutEntity(key, entity)
	db.addAof(utils.ToCmdLine2("set", args...))
	return reply.MakeOKReply()
}

// execSetNX SETNX k1 v1
func execSetNX(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := args[1]
	entity := &database.DataEntity{
		Data: val,
	}

	inserted := db.PutIfAbsent(key, entity)
	db.addAof(utils.ToCmdLine2("setnx", args...))
	return reply.MakeIntReply(int64(inserted))
}

// execGetSet GETSET k1 v1
func execGetSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := args[1]

	entity, exists := db.GetEntity(key)

	db.PutEntity(key, &database.DataEntity{
		Data: val,
	})

	if exists {
		db.addAof(utils.ToCmdLine2("getset", args...))
		return reply.MakeBulkReply(entity.Data.([]byte))
	}
	return reply.MakeNullBulkReply()
}

// execStrLen STRLEN
func execStrLen(db *DB, args [][]byte) resp.Reply {
	entity, exists := db.GetEntity(string(args[0]))
	if !exists {
		return reply.MakeNullBulkReply()
	}

	return reply.MakeIntReply(int64(len(entity.Data.([]byte))))
}
