package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jdnCreations/gator/internal/database"
	"github.com/lib/pq"
)

func handlerUnfollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feedurl>", cmd.Name)
	}

	err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		Name: user.Name,
		Url: cmd.Args[0],
	})
	if err != nil {
		return fmt.Errorf("could not delete follow: %v", err)
	}
	fmt.Printf("unfollowed feed: %s", cmd.Args[0])
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}


	feedfollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("could not get feed follows for userid: %v", user.ID)
	}

	if len(feedfollows) == 0 {
		 fmt.Println("not following any feeds")
		 return nil
	}

	fmt.Println(user.Name)
	for _, feed := range feedfollows {
		fmt.Println(feed.FeedName)
	}

	
	return nil
}

func handlerFollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feedurl>", cmd.Name)
	}

	url := cmd.Args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}
	
	ff, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			fmt.Println("You already follow this feed.")
			return nil
		}
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

	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// fmt.Printf("User not found for userID: %s\n", feed.UserID)
				continue
			}
			return err
		}

		fmt.Println("--------------------------")
		fmt.Printf("Feed: %s\n", feed.Name)
		fmt.Printf("Feed URL: %s\n", feed.Url)
		fmt.Printf("Feeds Adder: %s\n", user.Name)
	}
	return nil

}

func handlerAddFeed(s *state, cmd command, user database.User) error {
		if len(cmd.Args) != 2{
		return fmt.Errorf("usage: %s <feedname> <urloffeed>", cmd.Name)
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

	fmt.Printf("Added feed: %s, %s\n", feed.Name, feed.Url)
	return nil

}