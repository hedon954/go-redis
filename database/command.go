package database

import "strings"

// cmdTable holds all commands supported
var cmdTable = make(map[string]*command)

// command redis command wrapper
type command struct {
	executor ExecFun // command executor
	arity    int     // arg count
}

// RegisterCommand adds a command to cmdTable
func RegisterCommand(name string, executor ExecFun, arity int) {
	name = strings.ToLower(name)
	cmdTable[name] = &command{
		executor: executor,
		arity:    arity,
	}
}
