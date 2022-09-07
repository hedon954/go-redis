package reply

import (
	"bytes"
	"fmt"
	"strconv"

	"go-redis/interface/resp"
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
	return []byte(buildStringReply(b.Arg))
}

func MakeBulkReply(arg []byte) *BulkReply {
	return &BulkReply{
		Arg: arg,
	}
}

// MultiBulkReply represents multi messages redis replies to client
type MultiBulkReply struct {
	Args [][]byte
}

func (m *MultiBulkReply) ToBytes() []byte {
	argLen := len(m.Args)
	if argLen == 0 {
		return nullBulkReplyBytes
	}
	// SET key value
	// ->
	// *3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("*%d%s", argLen, CRLF))
	for i := 0; i < argLen; i++ { //$3\r\nSET\r\n
		if len(m.Args[i]) == 0 {
			buf.WriteString(string(nullBulkReplyBytes) + CRLF)
		} else {
			buf.WriteString(buildStringReply(m.Args[i]))
		}
	}
	return buf.Bytes()
}

func MakeMultiBulkReply(args [][]byte) *MultiBulkReply {
	return &MultiBulkReply{
		Args: args,
	}
}

// StatusReply replies a status
type StatusReply struct {
	Status string
}

func (s *StatusReply) ToBytes() []byte {
	return []byte("+" + s.Status + CRLF)
}

func MakeStatusReply(status string) *StatusReply {
	return &StatusReply{
		Status: status,
	}
}

// IntReply replies an int64 number
type IntReply struct {
	Code int64
}

func (i *IntReply) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(i.Code, 10) + CRLF)
}

func MakeIntReply(code int64) *IntReply {
	return &IntReply{
		Code: code,
	}
}

// ErrorReply defines error reply
type ErrorReply interface {
	Error() string
	ToBytes() []byte
}

// StandardErrReply defines an err reply can be customized
type StandardErrReply struct {
	Status string
}

func (s *StandardErrReply) Error() string {
	return s.Status
}

func (s *StandardErrReply) ToBytes() []byte {
	return []byte("-" + s.Status + CRLF)
}

func MakeStandardErrReply(status string) *StandardErrReply {
	return &StandardErrReply{
		Status: status,
	}
}

func buildStringReply(bs []byte) string {
	return fmt.Sprintf("$%d%s%s%s", len(bs), CRLF, bs, CRLF)
}

// IsErrReply checks the reply is error reply or not
func IsErrReply(reply resp.Reply) bool {
	return reply.ToBytes()[0] == '-'
}
