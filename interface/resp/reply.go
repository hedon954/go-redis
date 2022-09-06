package resp

// Reply defines the msg which the server passes to the client
type Reply interface {
	ToBytes() []byte
}
