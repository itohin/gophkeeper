package router

import (
	"fmt"
)

type Command interface{}

type Router struct {
	commands map[string]Command
}

func NewRouter(commands map[string]Command) *Router {
	return &Router{
		commands: commands,
	}
}

func (r *Router) GetCommand(key string) (Command, error) {
	cmd, ok := r.commands[key]
	if !ok {
		return nil, fmt.Errorf("unknown command %s", key)
	}
	return cmd, nil
}
