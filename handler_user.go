package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jdnCreations/gator/internal/database"
)

func handlerUsers(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting all users: %v", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUsername {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	err := s.db.DeleteAll(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("reset unsuccessful: %v", err)
	}

	fmt.Println("reset successful")
	os.Exit(0)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
	})

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return fmt.Errorf("user %s already exists", name)
		}	
		return fmt.Errorf("failed to create user: %v", err)
	}
	
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}

	fmt.Println("user was created:")
	fmt.Println(user)

	return nil

}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	user, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		fmt.Println("user does not exist")
		os.Exit(1)
		return err
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err) 
	}

	fmt.Println("User switched successfully!")
	return nil

}