package reply

import (
	"fmt"
)

var (
	nullBulkReplyBytes = []byte("$-1")
	CRLF               = "\r\n"
)

// BulkReply represents a msg redis replies to client
type BulkReply struct {
	Arg []byte
}

func (b *BulkReply) ToBytes() []byte {
	if len(b.Arg) == 0 {
		return nullBulkReplyBytes
	}
	// hedon -> $5\r\nhedon\r\n
	return []byte(fmt.Sprintf("$%d%s%s%s", len(b.Arg), CRLF, b.Arg, CRLF))
}

func MakeBulkReply(arg []byte) *BulkReply {
	return &BulkReply{
		Arg: arg,
	}
}

// ErrorReply defines error reply
type ErrorReply interface {
	Error() string
	ToBytes() []byte
}
