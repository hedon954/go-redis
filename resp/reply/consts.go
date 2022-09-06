package reply

/**
this file defines same fixed reply
*/

// PongReply pong
type PongReply struct {
}

var pongBytes = []byte("+PONG\r\n")

func (p *PongReply) ToBytes() []byte {
	return pongBytes
}

var thePongReply = new(PongReply)

func MakePongReply() *PongReply {
	return thePongReply
}

// OKReply OK
type OKReply struct {
}

var okBytes = []byte("+OK\r\n")

func (O OKReply) ToBytes() []byte {
	return okBytes
}

var theOKReplay = new(OKReply)

func MakeOKReply() *OKReply {
	return theOKReplay
}

// NullBulkReply empty reply
type NullBulkReply struct {
}

var nullBulkBytes = []byte("$-1\r\n")

func (n NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

var theNullBulkReply = new(NullBulkReply)

func MakeNullBulkReply() *NullBulkReply {
	return theNullBulkReply
}

// EmptyMultiBulkReply empty array
type EmptyMultiBulkReply struct {
}

var emptyMultiBulkBytes = []byte("*0\r\r")

func (e EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

var theEmptyMultiBulkReply = new(EmptyMultiBulkReply)

func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return theEmptyMultiBulkReply
}

// NoReply no reply
type NoReply struct {
}

var noBytes = []byte("")

func (n NoReply) ToBytes() []byte {
	return noBytes
}

var theNoReply = new(NoReply)

func MakeNoReply() *NoReply {
	return theNoReply
}
