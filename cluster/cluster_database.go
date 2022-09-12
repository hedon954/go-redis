package cluster

import (
	"go-redis/interface/database"
	"go-redis/interface/resp"
	"go-redis/lib/consistenthash"

	pool "github.com/jolestar/go-commons-pool/v2"
)

// Database cluster database
type Database struct {

	// self name
	self string

	// cluster nodes
	nodes []string

	// node selector
	peerPicker *consistenthash.NodeMap

	// connection pool
	peerConnection map[string]*pool.ObjectPool

	// inner db
	db database.Database
}

// MakeClusterDatabase creates a cluster database
func MakeClusterDatabase() *Database {

}

func (cdb *Database) Exec(client resp.Connection, args [][]byte) resp.Reply {
	//TODO implement me
	panic("implement me")
}

func (cdb *Database) Close() error {
	//TODO implement me
	panic("implement me")
}

func (cdb *Database) AfterClientClose(client resp.Connection) error {
	//TODO implement me
	panic("implement me")
}
