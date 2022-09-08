// Package database
// @description redis core service
package database

import "go-redis/interface/resp"

type CmdLine = [][]byte

type Database interface {
	Exec(client resp.Connection, args [][]byte) resp.Reply
	Close() error
	AfterClientClose(client resp.Connection) error
}

// DataEntity represents redis data structure
type DataEntity struct {
	Data interface{} // string, hash, list, set, sorted set
}
