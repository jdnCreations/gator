package main

import (
	"errors"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.registeredCommands[cmd.Name]
	if !exists {
		return errors.New("command does not exist")
	}

	return handler(s, cmd)
}