package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/alfredo-bertorelli/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %s", err)
	}
	req.Header.Set("User-Agent", "gator")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching feed: %s", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %s", err)
	}

	var rssFeed RSSFeed
	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling xml: %s", err)
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i, item := range rssFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		rssFeed.Channel.Item[i] = item
	}

	return &rssFeed, nil
}

func scrapeFeed(s *state) error {

	feedRow, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feed to fetch next: %w", err)
	}
	updatedFeed := database.UpdateLastFecthedAtParams{
		ID: feedRow.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: time.Now(),
	}
	err = s.db.UpdateLastFecthedAt(context.Background(), updatedFeed)
	if err != nil {
		return fmt.Errorf("couldn't mark feed as fetched: %w", err)
	}

	feed, err := fetchFeed(context.Background(), feedRow.Url)
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}

	newPosts := 0
	fmt.Printf("-- Feed: %s\n\n", feed.Channel.Title)
	for _, post := range feed.Channel.Item {
		postRow, err := storePost(s, post, feedRow.ID)
		if err != nil {
			return err
		} else if postRow != nil {
			newPosts++
			fmt.Printf("Post '%s' created!\n", postRow.Title)
		}
	}
	fmt.Printf("%d new posts stored\n\n", newPosts)
	return nil
}

func storePost(s *state, post RSSItem, feedId uuid.UUID) (*database.Post, error) {

	publishedAt := sql.NullTime{}
	if t, err := time.Parse(time.RFC1123Z, post.PubDate); err == nil {
		publishedAt = sql.NullTime{
			Time:  t,
			Valid: true,
		}
	}
	postParams := database.CreatePostParams{
		ID:    uuid.New(),
		Title: post.Title,
		Url:   post.Link,
		Description: sql.NullString{
			String: post.Description,
			Valid:  post.Description != "",
		},
		PublishedAt: publishedAt,
		FeedID:      feedId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	postRow, err := s.db.CreatePost(context.Background(), postParams)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("couldn't create post: %w", err)
	}

	return &postRow, nil
}
