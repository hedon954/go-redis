// Package parse
// @description provides the capabilities to parse the msg client sent to redis server
package parse

import (
	"bufio"
	"errors"
	"fmt"
	"go-redis/interface/resp"
	"go-redis/lib/logger"
	"go-redis/resp/reply"
	"io"
	"runtime/debug"
	"strconv"
	"strings"
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
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()

	bufReader := bufio.NewReader(reader)
	var state readState
	var err error
	var msg []byte

	for true {
		var ioErr bool
		msg, ioErr, err = readLine(bufReader, &state)
		if err != nil {
			// finished server current client while io err occurs
			if ioErr {
				ch <- &Payload{
					Err: err,
				}
				close(ch)
				return
			}
			// protocol error
			ch <- &Payload{
				Err: err,
			}
			state = readState{}
			continue
		}

		if len(msg) <= 0 {
			continue
		}

		if !state.readingMultiLine {
			// reading single line mode, transfer to reading multi lines mode
			if msg[0] == '*' {
				err = parseMultiBulkHeader(msg, &state)
				if err != nil {
					ch <- &Payload{
						Err: errors.New("protocol error: " + string(msg)),
					}
					state = readState{}
					continue
				}
				if state.expectedArgsCount == 0 {
					ch <- &Payload{
						Data: reply.MakeEmptyMultiBulkReply(),
					}
					state = readState{}
					continue
				}
			} else if msg[0] == '$' { // $5\r\nhedon\r\n
				err = parseBulkHeader(msg, &state)
				if err != nil {
					ch <- &Payload{
						Err: errors.New("protocol error: " + string(msg)),
					}
					state = readState{}
					continue
				}
				if state.bulkLen == -1 { // $-1\r\n
					ch <- &Payload{
						Data: reply.MakeNullBulkReply(),
					}
					state = readState{}
					continue
				}
			} else { // + - :
				result, err := parseSingleLineReply(msg)
				ch <- &Payload{
					Data: result,
					Err:  err,
				}
				state = readState{}
				continue
			}
		} else {
			// reading multi line mode, reads body
			err = readBody(msg, &state)
			if err != nil {
				ch <- &Payload{
					Err: errors.New("protocol error: " + string(msg)),
				}
				state = readState{}
				continue
			}

			if state.finished() {
				var result resp.Reply
				if state.msgType == '*' {
					result = reply.MakeMultiBulkReply(state.args)
				} else if state.msgType == '$' {
					result = reply.MakeBulkReply(state.args[0])
				}
				ch <- &Payload{
					Data: result,
				}
				state = readState{}
				continue
			}
		}

		s.ParseStream(reader)
	}
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

// parseSingleLineReply parses single line reply, gets inner message
// example:
// 	+OK\r\n
//  -Err\r\n
//  :5\r\n
func parseSingleLineReply(line []byte) (resp.Reply, error) {
	str := string(line)
	str = strings.TrimSuffix(str, "\r\n")
	if len(str) < 1 {
		return nil, fmt.Errorf("protocol error: %s", string(line))
	}
	var res resp.Reply
	switch str[0] {
	case '+':
		res = reply.MakeStatusReply(str[1:])
	case '-':
		res = reply.MakeStandardErrReply(str[1:])
	case ':':
		i, err := strconv.ParseInt(str[1:], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("protocol error: %s", string(line))
		}
		res = reply.MakeIntReply(i)
	}
	return res, nil
}

// readBody reads body from command message
// example:
// 	hedon\r\n
// 	$3\r\nSET\r\n$3\r\nkey\r\n$9\r\nval\r\ue\r\n
func readBody(line []byte, state *readState) error {
	if len(line) < 3 {
		return fmt.Errorf("protocol error: %s", string(line))
	}
	// remove \r\n
	line = line[0 : len(line)-2]
	var err error

	// $3
	if line[0] == '$' {
		// $3 -> 3
		state.bulkLen, err = strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return fmt.Errorf("protocol error: %s", string(line))
		}
		if state.bulkLen <= 0 { //$0\r\n
			state.args = append(state.args, []byte{})
			state.bulkLen = 0
		}
	} else {
		// SET\r\n
		state.args = append(state.args, line)
	}

	return nil
}
