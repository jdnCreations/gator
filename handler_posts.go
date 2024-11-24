package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jdnCreations/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %s <optional limit>", cmd.Name)
	}
	
	limit := 2

	if len(cmd.Args) == 1 {
		parsedLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return err
		}

		limit = parsedLimit
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: int32(limit),
	})
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Println("---------------------")
		fmt.Printf("Title: %s\n",post.Title)
		fmt.Printf("Url: %s\n",post.Url)
		fmt.Printf("Description: %s\n",post.Description)
		fmt.Printf("Publish at: %s\n", post.PublishedAt)
		feed, err := s.db.GetFeedById(context.Background(), post.FeedID)
		if err != nil {
			return err
		}
		fmt.Printf("Feed: %s\n", feed.Url)
	}
	
	if len(posts) == 0 {
		fmt.Println("No posts found!")
	}
	return nil;
}