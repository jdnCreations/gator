package main

import (
	"context"
	"errors"

	"github.com/jdnCreations/gator/internal/database"
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

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func (*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
		if err != nil {
			return err
		}
		err = handler(s, cmd, user)
		if err != nil {
			return err
		}
		return nil
	}
}