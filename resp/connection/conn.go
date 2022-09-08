package connection

import (
	"net"
	"sync"
	"time"

	"go-redis/lib/sync/wait"
)

// Connection redis client connection
type Connection struct {
	conn         net.Conn   // client tcp connection
	waitingReply wait.Wait  // waiting until reply finished
	mu           sync.Mutex // lock while handler sending response
	selectedDB   int        // selected redis db
}

// NewConn creates a new connection
func NewConn(conn net.Conn) *Connection {
	return &Connection{
		conn: conn,
	}
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// Write sends msg to client
func (c *Connection) Write(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}
	c.mu.Lock()
	c.waitingReply.Add(1)
	defer func() {
		c.waitingReply.Done()
		c.mu.Unlock()
	}()

	_, err := c.conn.Write(bytes)
	return err
}

func (c *Connection) GetDBIndex() int {
	return c.selectedDB
}

func (c *Connection) SelectDB(dbNum int) {
	c.selectedDB = dbNum
}

func (c *Connection) Close() error {
	c.waitingReply.WaitWithTimeout(10 * time.Second)
	_ = c.conn.Close()
	return nil
}
