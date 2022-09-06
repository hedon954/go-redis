package resp

// Connection defines a redis connection
type Connection interface {
	Write([]byte) error // writes data to client
	GetDBIndex() int    // redis has multi databases
	SelectDB(int)       // select redis database
}
