package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jdnCreations/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
  if len(cmd.Args) != 1 {
    return fmt.Errorf("usage: %s <time between reqs>", cmd.Name)
  }

  timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
  if err != nil {
    return fmt.Errorf("invalid duration: %w", err) 
  }

	fmt.Printf("collecting feeds every %d seconds\n", timeBetweenReqs)
  ticker := time.NewTicker(timeBetweenReqs)

  for ; ; <-ticker.C {
    scrapeFeeds(s)
  }
}


func scrapeFeeds(s *state) {
  feed, err := s.db.GetNextFeedToFetch(context.Background())
  if err != nil {
		fmt.Println("Couldn't get next feeds to fetch", err)
    return
  }
	fmt.Println("Found a feed to fetch!")
	scrapeFeed(s.db, feed)
}


func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
  if err != nil {
		fmt.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
    return
  }
  feedData, err := fetchFeed(context.Background(), feed.Url)
  if err != nil {
		fmt.Printf("Couldn't collect feed: %s: %v", feed.Name, err)
		return
  }
	for _, item := range feedData.Channel.Item {
			publishedAt := sql.NullTime{}
			if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
				publishedAt = sql.NullTime{
					Time: t,
					Valid: true,
				}
			}

			_, err = db.CreatePost(context.Background(), database.CreatePostParams{
				ID: uuid.New(),
				Title: item.Title,
				Url: item.Link,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
				Description: sql.NullString{
					String: item.Description,
					Valid: true,
				},
				PublishedAt: publishedAt,
				FeedID: feed.ID,
			})
			if err != nil {
				if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
			}
	}

	fmt.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}