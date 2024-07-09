package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Evilcmd/GO-Blog-Aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiConfig *apiConfigDefn) createFeeds(res http.ResponseWriter, req *http.Request) {
	type createFeedStructDefn struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(req.Body)
	createFeedRes := createFeedStructDefn{}
	err := decoder.Decode(&createFeedRes)
	if err != nil {
		respondWithError(res, 400, fmt.Sprintf("problem reading the req body %v", err.Error()))
		return
	}

	feed, err := apiConfig.DB.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   sql.NullString{String: createFeedRes.Name, Valid: true},
		Url:    sql.NullString{String: createFeedRes.URL, Valid: true},
		UserID: uuid.NullUUID{UUID: apiConfig.user.ID, Valid: true},
	})

	if err != nil {
		respondWithError(res, 400, fmt.Sprintf("error while creating feed: %v", err.Error()))
		return
	}

	feedFollowDB, err := apiConfig.helperFuncTCreateFeedFollow(feed.ID)
	if err != nil {
		respondWithError(res, 200, fmt.Sprintf("error while adding to database: %v", err.Error()))
		return
	}

	type ReturnFeed struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		URL       string    `json:"url"`
		UserID    string    `json:"user_id"`
	}

	type ActualReturn struct {
		Feed       ReturnFeed              `json:"feed"`
		FeedFollow FeedFollowresponseStrct `json:"feed_follow"`
	}

	x := ActualReturn{
		Feed: ReturnFeed{
			ID:        feed.ID.String(),
			CreatedAt: apiConfig.user.CreatedAt,
			UpdatedAt: apiConfig.user.UpdatedAt,
			Name:      feed.Name.String,
			URL:       feed.Url.String,
			UserID:    apiConfig.user.ID.String(),
		},
		FeedFollow: feedFollowDB,
	}

	respondWithJson(res, 200, x)
}

func (apiConfig *apiConfigDefn) GetAllFeeds(res http.ResponseWriter, req *http.Request) {
	feeds, err := apiConfig.DB.GetAllFeeds(context.Background())
	if err != nil {
		respondWithError(res, 400, fmt.Sprintf("Error while retreiving data: %v", err.Error()))
		return
	}
	type ReturnFeeds struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Url    string `json:"url"`
		UserId string `json:"user_id"`
	}

	returnVal := make([]ReturnFeeds, 0, len(feeds))
	for _, v := range feeds {
		returnVal = append(returnVal, ReturnFeeds{
			ID:     v.ID.String(),
			Name:   v.Name.String,
			Url:    v.Url.String,
			UserId: v.UserID.UUID.String(),
		})
	}

	respondWithJson(res, 200, returnVal)
}
