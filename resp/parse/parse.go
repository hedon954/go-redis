// Package parse
// @description provides the capabilities to parse the msg client sent to redis server
package parse

import (
	"bufio"
	"fmt"
	"go-redis/interface/resp"
	"io"
	"strconv"
)

// Payload represents the client command
type Payload struct {

	// the format of the data transmitted
	// by the redis client and the server
	// is the same
	Data resp.Reply
	Err  error
}

// readState represents the status of the parser
type readState struct {

	// the parser is reading single line or multi lines
	// false -> single line
	// true  -> multi line
	readingMultiLine bool

	// the number of current command's args expected
	expectedArgsCount int

	// the type of message
	msgType byte

	// the command sent by client
	args [][]byte

	// the length of data bulk
	bulkLen int64
}

// finished checks the parse operation is finished or not
func (s *readState) finished() bool {
	return s.expectedArgsCount > 0 && len(s.args) == s.expectedArgsCount
}

// ParseStream parses message stream async
func (s *readState) ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)

	// async parse
	go s.parse0(reader, ch)
	return ch
}

// parse0 parses message
func (s *readState) parse0(reader io.Reader, ch chan<- *Payload) {

}

// readLine reads a line from bufReader
func readLine(bufReader *bufio.Reader, state *readState) (line []byte, ioErr bool, err error) {
	// *3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$9\r\nval\r\ue\r\n
	// 	  ↓
	// *3\r\n
	// 				↓
	//		  $3\r\nSET\r\n
	// 							 ↓
	//		  				$3\r\nkey\r\n
	//												↓
	//										$9\r\nval\r\nue\r\n

	// case 1. without '${bytes number}' -> split with '\r\n'
	if state.bulkLen == 0 {
		line, err = bufReader.ReadBytes('\n')
		if err != nil {
			return nil, true, err
		}
		if len(line) < 2 || line[len(line)-2] != '\r' {
			return nil, false, fmt.Errorf("protocol error: %s", string(line))
		}
		// return *{number}\r\n
	}

	// case 2. with '${bytes number}' -> split by 'bytes number'
	if state.bulkLen > 0 {
		line = make([]byte, state.bulkLen+2) // 2: \r\n
		_, err = io.ReadFull(bufReader, line)
		if err != nil {
			return nil, true, err
		}
		if len(line) < 2 || line[len(line)-2] != '\r' || line[len(line)-1] != '\n' {
			return nil, false, fmt.Errorf("protocol error: %s", string(line))
		}

		state.bulkLen = 0
		// return {msg}\r\n
	}

	return line, false, nil
}

// parseMultiBulkHeader sets readState according to '*{line}\r\n'
func parseMultiBulkHeader(line []byte, state *readState) error {
	var err error
	var expectedLine uint64

	// *3\r\n  -> 3
	// *300\r\n -> 300
	if len(line) < 3 || line[0] != '*' {
		return fmt.Errorf("protocol error, cannot get expectedLine: %s", line)
	}
	expectedLine, err = strconv.ParseUint(string(line[1:len(line)-2]), 10, 32)
	if err != nil {
		return fmt.Errorf("protocol error, cannot get expectedLine: %s, err is: %v", line, err)
	}

	if expectedLine == 0 {
		state.expectedArgsCount = 0
		return nil
	}

	state.expectedArgsCount = int(expectedLine)  // args count
	state.msgType = line[0]                      // reading array
	state.readingMultiLine = true                // reading multi line
	state.args = make([][]byte, 0, expectedLine) // initial args array
	return nil
}

// parseBulkHeader sets readState according to '${number}\r\n'
func parseBulkHeader(line []byte, state *readState) error {
	var err error
	var bulkLen uint64

	if len(line) < 2 || line[0] != '$' {
		return fmt.Errorf("protocol error: %s", string(line))
	}

	// $300\r\n -> 300
	bulkLen, err = strconv.ParseUint(string(line[1:len(line)-2]), 10, 32)
	if err != nil {
		return fmt.Errorf("protocol error: %s", string(line))
	}

	// null bulk
	if bulkLen <= 0 {
		return nil
	}

	state.msgType = line[0]
	state.readingMultiLine = true
	state.expectedArgsCount = 1
	state.args = make([][]byte, 0, 1)
	return nil
}
