package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Evilcmd/GO-Blog-Aggregator/internal/database"
	"github.com/google/uuid"
)

type Feed struct {
	Title       string `xml:"title"`
	URL         string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

type RssFeed struct {
	Name        string `xml:"channel>title"`
	URL         string `xml:"channel>link"`
	Description string `xml:"channel>description"`
	Feeds       []Feed `xml:"channel>item"`
}

func getRssFeed(db *database.Queries, url string, wg *sync.WaitGroup, feedId uuid.UUID) {
	defer wg.Done()

	httpClient := http.Client{Timeout: time.Second * 10}
	httpres, err := httpClient.Get(url)
	if err != nil {
		log.Printf("error getting the url: %v\n", err.Error())
		return
	}
	defer httpres.Body.Close()

	newRssFeed := RssFeed{}
	decoder := xml.NewDecoder(httpres.Body)
	err = decoder.Decode(&newRssFeed)
	if err != nil {
		log.Printf("error decoding the xml: %v\n", err.Error())
		return
	}

	for _, feed := range newRssFeed.Feeds {
		desc := sql.NullString{String: "", Valid: false}
		if len(feed.Description) > 0 {
			desc.String = newRssFeed.Description
			desc.Valid = true
		}
		t, err := time.Parse(time.RFC1123Z, feed.PubDate)
		pubAt := sql.NullTime{Time: time.Time{}, Valid: false}
		if err == nil {
			pubAt.Time = t
			pubAt.Valid = true
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       feed.Title,
			Url:         feed.URL,
			Description: desc,
			PublishedAt: pubAt,
			FeedID:      feedId,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("error in storing the feed: %v in database: %v\n", feed.URL, err.Error())
		}
	}

	// log.Println(newRssFeed)
	log.Println(url, len(newRssFeed.Feeds))
}

func rssWorker(db *database.Queries, concurrencyCount int, duration time.Duration) {

	timeTicker := time.NewTicker(duration)
	for ; ; <-timeTicker.C {
		log.Println("_______Start fetching from rss_______")
		feedsToFetch, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrencyCount))
		if err != nil {
			log.Println("error fetching from database")
			continue
		}
		wg := sync.WaitGroup{}
		for _, v := range feedsToFetch {
			wg.Add(1)
			go getRssFeed(db, v.Url.String, &wg, v.ID)
		}
		wg.Wait()
		log.Println("_______Done fetching from rss_______")
	}
}
