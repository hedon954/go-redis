package database

import (
	"strconv"
	"strings"

	"go-redis/aof"
	"go-redis/config"
	"go-redis/interface/resp"
	"go-redis/lib/logger"
	"go-redis/resp/reply"
)

// Database represent a redis
type Database struct {
	dbSet      []*DB
	aofHandler *aof.Handler
}

// NewDatabase initials a redis
func NewDatabase() *Database {
	// create dbs
	database := &Database{}
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}

	database.dbSet = make([]*DB, config.Properties.Databases)
	for i := range database.dbSet {
		db := makeDB()
		db.index = i
		database.dbSet[i] = db
	}

	// initial aof
	if config.Properties.AppendOnly {
		aofHandler, err := aof.NewAofHandler(database)
		if err != nil {
			panic(err)
		}
		database.aofHandler = aofHandler
		for _, db := range database.dbSet {
			db.addAof = func(line CmdLine) {
				database.aofHandler.AddAof(db.index, aof.CmdLine(line))
			}
		}

	}
	return database
}

// Exec executes command sent by client
func (database *Database) Exec(client resp.Connection, args [][]byte) resp.Reply {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()

	cmdName := strings.ToLower(string(args[0]))
	if cmdName == "select" {
		if len(args) != 2 {
			return reply.MakeArgNumErrReply("select")
		}

		return execSelect(client, database, args[1:])
	}

	db := database.dbSet[client.GetDBIndex()]
	return db.Exec(client, args)
}

func (database *Database) Close() error {
	logger.Info("database shutting down")
	return nil
}

func (database *Database) AfterClientClose(client resp.Connection) error {
	logger.Info("client shutting down")
	return nil
}

// execSelect selects a db
// e.g. select 1
func execSelect(c resp.Connection, database *Database, args [][]byte) resp.Reply {
	index := string(args[0])
	i, err := strconv.Atoi(index)
	if err != nil {
		return reply.MakeStandardErrReply("ERR invalid DB index")
	}

	if i >= len(database.dbSet) {
		return reply.MakeStandardErrReply("ERR DB index is out of range")
	}

	c.SelectDB(i)

	return reply.MakeOKReply()
}
