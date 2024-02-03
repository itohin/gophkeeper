package router

import "errors"

type Command func() (string, error)

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
		return nil, errors.New("unknown command")
	}
	return cmd, nil
}
