package reply

import "fmt"

// UnknownErrReply means unknown error occurs
type UnknownErrReply struct {
}

var unknownErrBytes = []byte("-Err unknown\r\n")

func (u *UnknownErrReply) Error() string {
	return "Err unknown"
}

func (u *UnknownErrReply) ToBytes() []byte {
	return unknownErrBytes
}

var theUnknownErrReply = new(UnknownErrReply)

func MakeUnknownErrReplay() *UnknownErrReply {
	return theUnknownErrReply
}

// ArgNumErrReply means the number of args sent by client is faulted
type ArgNumErrReply struct {
	Cmd string
}

func (a *ArgNumErrReply) Error() string {
	return fmt.Sprintf("ERR wrong number of arguments for '%s' command", a.Cmd)
}

func (a *ArgNumErrReply) ToBytes() []byte {
	return []byte(fmt.Sprintf("-ERR wrong number of arguments for '%s' command\r\n", a.Cmd))
}

func MakeArgNumErrReply(cmd string) *ArgNumErrReply {
	return &ArgNumErrReply{
		Cmd: cmd,
	}
}

// SyntaxErrReply means the command sent by client has syntax error
type SyntaxErrReply struct {
}

var syntaxErrBytes = []byte("-Err syntax error\r\n")

func (s *SyntaxErrReply) Error() string {
	return "Err syntax error"
}

var theSyntaxErrReply = new(SyntaxErrReply)

func (s *SyntaxErrReply) ToBytes() []byte {
	return syntaxErrBytes
}

func MakeSyntaxErrReply() *SyntaxErrReply {
	return theSyntaxErrReply
}

// WrongTypeErrReply means against a key holding the wrong kind of value
type WrongTypeErrReply struct {
}

var wrongTypeErrBytes = []byte("-Err Wrong type operation against a key holding the wrong kind of value\r\n")

func (w *WrongTypeErrReply) Error() string {
	return "Err Wrong type operation against a key holding the wrong kind of value"
}

func (w *WrongTypeErrReply) ToBytes() []byte {
	return wrongTypeErrBytes
}

var theWrongTypeErrBytes = new(WrongTypeErrReply)

func MakeWrongTypeErrReply() *WrongTypeErrReply {
	return theWrongTypeErrBytes
}

// ProtocolErrReply means the msg sent by client is compliant with RESP
type ProtocolErrReply struct {
	Msg string
}

func (p *ProtocolErrReply) Error() string {
	return fmt.Sprintf("ERR Protocol error: '%s'", p.Msg)
}

func (p *ProtocolErrReply) ToBytes() []byte {
	return []byte(fmt.Sprintf("-ERR Protocol error: '%s' \r\n", p.Msg))
}

func MakeProtocolErrReply(msg string) *ProtocolErrReply {
	return &ProtocolErrReply{
		Msg: msg,
	}
}
