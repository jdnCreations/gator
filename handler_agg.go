package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jdnCreations/gator/internal/database"
	"github.com/lib/pq"
)

func handlerAgg(s *state, cmd command) error {
  if len(cmd.Args) != 1 {
    return fmt.Errorf("usage: %s <time between reqs>", cmd.Name)
  }

  timeBetweenReqs, err := strconv.Atoi(cmd.Args[0])
  if err != nil {
    return err
  }

  ticker := time.NewTicker(time.Duration(timeBetweenReqs) * time.Second)

  for ; ; <-ticker.C {
    fmt.Printf("collecting feeds every %d seconds\n", timeBetweenReqs)
    scrapeFeeds(s)
  }
}

func scrapeFeeds(s *state) error {
  next, err := s.db.GetNextFeedToFetch(context.Background())
  if err != nil {
    return err
  }
	fmt.Printf("Fetching feed: %s\n", next.Url)
  err = s.db.MarkFeedFetched(context.Background(), next.ID)
  if err != nil {
    return err
  }
  rss, err := fetchFeed(context.Background(), next.Url)
  if err != nil {
    return err
  }
	
	fmt.Printf("found %d items in feed\n", len(rss.Channel.Item))

  for i := range rss.Channel.Item {
		layout := "Mon, 02 Jan 2006 15:04:05 -0700"

		fmt.Printf("Attemping to parse date: %s\n", rss.Channel.Item[i].PubDate)
		pubDate, err := time.Parse(layout, rss.Channel.Item[i].PubDate)
		if err != nil {
			fmt.Printf("Error parsing date: %v\n", err) 
			return err
		}

		fmt.Printf("Attemping to create post: \n")
		fmt.Printf("  Title: %s\n", rss.Channel.Item[i].Title)
		fmt.Printf("  URL: %s\n", rss.Channel.Item[i].Link)
		fmt.Printf("  Description: %s\n", rss.Channel.Item[i].Description)
		fmt.Printf("  Published At: %v\n", rss.Channel.Item[i].PubDate)
		fmt.Printf("  Feed ID: %s\n", next.ID)

		post, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			Title: rss.Channel.Item[i].Title,
			Url: rss.Channel.Item[i].Link,
			Description: rss.Channel.Item[i].Description,
			PublishedAt: pubDate,
			FeedID: next.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "unique violation") {
				fmt.Printf("Skipping duplicate post: %s\n", rss.Channel.Item[i].Title)
				continue
			} 

			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
				fmt.Printf("Skipping duplicate post: %s\n", rss.Channel.Item[i].Title)
				continue
			} else {
				fmt.Printf("Error creating post: %v\n", err)
			}

			return err
		}

		fmt.Printf("created post: %v\n", post.Title)
  }

  return nil
}