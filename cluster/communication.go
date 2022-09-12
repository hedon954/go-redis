package cluster

import (
	"context"
	"fmt"
	"go-redis/interface/resp"
	"go-redis/lib/utils"
	"go-redis/resp/reply"
	"strconv"

	"go-redis/resp/client"
)

// getPeerClient gets a connection client from pool
func (cdb *Database) getPeerClient(peer string) (*client.Client, error) {
	pool, ok := cdb.peerConnection[peer]
	if !ok {
		return nil, fmt.Errorf("connection not found")
	}

	object, err := pool.BorrowObject(context.Background())
	if err != nil {
		return nil, err
	}

	c, ok := object.(*client.Client)
	if !ok {
		return nil, fmt.Errorf("type mismatch")
	}
	return c, nil
}

// returnPeerClient sends connection client back to pool
func (cdb *Database) returnPeerClient(peer string, c *client.Client) error {
	pool, ok := cdb.peerConnection[peer]
	if !ok {
		return fmt.Errorf("connection not found")
	}

	return pool.ReturnObject(context.Background(), c)
}

// relay forwards the request to specific redis node
func (cdb *Database) relay(peer string, c resp.Connection, args [][]byte) resp.Reply {

	// peer is current node
	if peer == cdb.self {
		return cdb.db.Exec(c, args)
	}

	// call peer node
	peerClient, err := cdb.getPeerClient(peer)
	if err != nil {
		return reply.MakeStandardErrReply(err.Error())
	}
	defer func() {
		_ = cdb.returnPeerClient(peer, peerClient)
	}()

	// select db
	peerClient.Send(utils.ToCmdLine("SELECT", strconv.Itoa(c.GetDBIndex())))

	// send command
	return peerClient.Send(args)
}

// broadcast commands to cluster nodes
func (cdb *Database) broadcast(c resp.Connection, args [][]byte) map[string]resp.Reply {
	res := make(map[string]resp.Reply)

	for _, node := range cdb.nodes {
		result := cdb.relay(node, c, args)
		res[node] = result
	}
	return res
}
