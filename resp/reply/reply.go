package reply

// ErrorReply defines error reply
type ErrorReply interface {
	Error() string
	ToBytes() []byte
}
