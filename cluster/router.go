package cluster

import (
	"go-redis/interface/resp"
)

// CmdFunc represents a command executor
type CmdFunc func(cdb *Database, c resp.Connection, cmdArgs [][]byte) resp.Reply

// makeRouter creates a router
func makeRouter() map[string]CmdFunc {
	m := make(map[string]CmdFunc)

	// need relay
	m["rename"] = renameFunc   // rename src dst
	m["renamenx"] = renameFunc // renamenx src dst
	m["exists"] = defaultFunc  // exists k1
	m["get"] = defaultFunc     // get k1
	m["set"] = defaultFunc     // set k1 v1
	m["setnx"] = defaultFunc   // setnx k1 v1
	m["getset"] = defaultFunc  // getset k1 v1
	m["type"] = defaultFunc    // type k1

	// need not relay
	m["ping"] = pingFunc     // ping
	m["select"] = selectFunc // select 1

	// need to broadcast
	m["flushdb"] = flushdbFunc // flushdb
	m["delete"] = delFunc      // delete k1 k2 k3 ...

	return m
}

// defaultFunc represent a default CmdFunc, used while need to relay command
// GET K1
// SET K1 V1
// EXISTS K1
func defaultFunc(cdb *Database, c resp.Connection, cmdArgs [][]byte) resp.Reply {

	// select node
	key := string(cmdArgs[1])
	node := cdb.peerPicker.PickNode(key)

	// relay
	return cdb.relay(node, c, cmdArgs)
}
