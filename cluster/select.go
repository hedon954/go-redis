package cluster

import "go-redis/interface/resp"

func selectFunc(cdb *Database, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	return cdb.db.Exec(c, cmdArgs)
}
