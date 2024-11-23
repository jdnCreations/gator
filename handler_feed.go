package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jdnCreations/gator/internal/database"
)

func handlerFollowing(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	username := s.cfg.CurrentUsername
	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("error getting user: %v", err)
	}

	feedfollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("could not get feed follows for userid: %v", user.ID)
	}

	if len(feedfollows) == 0 {
		 fmt.Println("not following any feeds")
		 return nil
	}

	fmt.Println(username)
	for _, feed := range feedfollows {
		fmt.Println(feed.FeedName)
	}

	
	return nil
}

func handlerFollowFeed(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feedurl>", cmd.Name)
	}

	url := cmd.Args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
	if err != nil {
		return err
	}
	
	ff, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Println(ff.FeedName)
	fmt.Println(ff.UserName)

	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}
	
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	fmt.Println(feeds)

	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				fmt.Printf("User not found for userID: %s\n", feed.UserID)
				continue
			}
			return err
		}

		fmt.Println(feed.Name)
		fmt.Println(feed.Url)
		fmt.Println(user.Name)
	}
	return nil

}

func handlerAddFeed(s *state, cmd command) error {
		if len(cmd.Args) != 2{
		return fmt.Errorf("usage: %s <feedname> <urloffeed>", cmd.Name)
	}

	userName := s.cfg.CurrentUsername
	user, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		return err
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
		Url: url,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("error getting all users: %v", err)
	}
	
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating follow feed when creating a new feed: %v", err)
	}

	fmt.Println(feed.ID)
	fmt.Println(feed.CreatedAt)
	fmt.Println(feed.UpdatedAt)
	fmt.Println(feed.Name)
	fmt.Println(feed.Url)
	fmt.Println(feed.UserID)

	return nil

}