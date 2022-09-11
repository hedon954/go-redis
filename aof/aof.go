package aof

import (
	"os"
	"strconv"

	"go-redis/config"
	"go-redis/lib/logger"
	"go-redis/lib/utils"
	"go-redis/resp/reply"

	databaseface "go-redis/interface/database"
)

// CmdLine is alias for [][]byte, represents a command line
type CmdLine [][]byte

const (
	aofQueueSize = 1 << 16
)

type payload struct {
	cmdLine CmdLine
	dbIndex int
}

// Handler receives messages from channel and write to AOF file
type Handler struct {
	db          databaseface.Database
	aofChan     chan *payload
	aofFile     *os.File
	aofFilename string
	currentDB   int
}

// NewAofHandler creates a new aof Handler
func NewAofHandler(db databaseface.Database) (*Handler, error) {
	handler := &Handler{}
	handler.aofFilename = config.Properties.AppendFilename
	handler.db = db

	handler.loadAof()

	aofFile, err := os.OpenFile(handler.aofFilename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	handler.aofFile = aofFile

	handler.aofChan = make(chan *payload, aofQueueSize)
	go func() {
		handler.handleAof()
	}()

	return handler, nil
}

// AddAof send command to aof goroutine through channel
func (handler *Handler) AddAof(dbIndex int, cmdLine CmdLine) {
	if config.Properties.AppendOnly && handler.aofChan != nil {
		handler.aofChan <- &payload{
			cmdLine: cmdLine,
			dbIndex: dbIndex,
		}
	}
}

// handleAof listens aof channel and write into file
func (handler *Handler) handleAof() {
	handler.currentDB = 0
	for p := range handler.aofChan {
		if p.dbIndex != handler.currentDB {
			// select db
			data := reply.MakeMultiBulkReply(utils.ToCmdLine("SELECT", strconv.Itoa(p.dbIndex))).ToBytes()
			_, err := handler.aofFile.Write(data)
			if err != nil {
				logger.Warn(err)
				continue // skip this command
			}
			handler.currentDB = p.dbIndex
		}
		data := reply.MakeMultiBulkReply(p.cmdLine).ToBytes()
		_, err := handler.aofFile.Write(data)
		if err != nil {
			logger.Warn(err)
		}
	}
}

// loadAof replays aof file when redis starts
func (handler *Handler) loadAof() {

}
