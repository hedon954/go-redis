package cluster

import (
	"context"
	pool "github.com/jolestar/go-commons-pool/v2"
)

// connectionFactory provides the operations to control connection pool
// implements pool.PooledObjectFactory interface
type connectionFactory struct {
	Peer string
}

func (c connectionFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	//TODO implement me
	panic("implement me")
}

func (c connectionFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	//TODO implement me
	panic("implement me")
}

func (c connectionFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	//TODO implement me
	panic("implement me")
}

func (c connectionFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	//TODO implement me
	panic("implement me")
}

func (c connectionFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	//TODO implement me
	panic("implement me")
}
