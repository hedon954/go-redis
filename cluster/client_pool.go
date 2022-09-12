package cluster

import (
	"context"
	"fmt"

	"go-redis/resp/client"

	pool "github.com/jolestar/go-commons-pool/v2"
)

// connectionFactory provides the operations to control connection pool
// implements pool.PooledObjectFactory interface
type connectionFactory struct {

	// node address
	Peer string
}

func (cf connectionFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	redisClient, err := client.MakeClient(cf.Peer)
	if err != nil {
		return nil, err
	}

	redisClient.Start()
	return pool.NewPooledObject(redisClient), nil
}

func (cf connectionFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	c, ok := object.Object.(*client.Client)
	if !ok {
		return fmt.Errorf("type mismatch")
	}

	c.Close()
	return nil
}

func (cf connectionFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	return true
}

func (cf connectionFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}

func (cf connectionFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}
