package cluster

import (
	"context"
	"runtime/debug"
	"strings"

	"go-redis/config"
	"go-redis/interface/database"
	"go-redis/interface/resp"
	"go-redis/lib/consistenthash"
	"go-redis/lib/logger"
	"go-redis/resp/reply"

	database2 "go-redis/database"

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
	cluster := &Database{
		self:           config.Properties.Self,
		db:             database2.NewStandaloneDatabase(),
		peerPicker:     consistenthash.NewNodeMap(nil),
		peerConnection: make(map[string]*pool.ObjectPool),
	}

	// nodes
	nodes := make([]string, 0, len(config.Properties.Peers)+1)
	for _, peer := range config.Properties.Peers {
		nodes = append(nodes, peer)
	}
	nodes = append(nodes, config.Properties.Self)
	cluster.nodes = nodes

	// peerPicker
	cluster.peerPicker.AddNode(nodes...)

	// node pool
	ctx := context.Background()
	for _, peer := range config.Properties.Peers {
		cluster.peerConnection[peer] = pool.NewObjectPoolWithDefaultConfig(ctx, &connectionFactory{
			Peer: peer,
		})
	}

	return cluster
}

var router = makeRouter()

func (cdb *Database) Exec(c resp.Connection, args [][]byte) (result resp.Reply) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			logger.Error(err)
			result = reply.MakeUnknownErrReplay()
		}
	}()

	// get command func
	cmdFunc, ok := router[strings.ToLower(string(args[0]))]
	if !ok {
		return reply.MakeStandardErrReply("ERR command unsupported: " + string(args[0]))
	}

	result = cmdFunc(cdb, c, args)
	return
}

func (cdb *Database) Close() error {
	return cdb.db.Close()
}

func (cdb *Database) AfterClientClose(c resp.Connection) error {
	return cdb.db.AfterClientClose(c)
}
