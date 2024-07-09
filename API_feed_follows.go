package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Evilcmd/GO-Blog-Aggregator/internal/database"
	"github.com/google/uuid"
)

type FeedFollowresponseStrct struct {
	ID        string    `json:"id"`
	FeedID    string    `json:"feed_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (apiConfig *apiConfigDefn) helperFuncTCreateFeedFollow(FeedID uuid.UUID) (FeedFollowresponseStrct, error) {
	feedFollowDB, err := apiConfig.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    apiConfig.user.ID,
		FeedsID:   FeedID,
	})
	if err != nil {
		return FeedFollowresponseStrct{}, nil
	}
	x := FeedFollowresponseStrct{
		ID:        feedFollowDB.ID.String(),
		FeedID:    feedFollowDB.FeedsID.String(),
		UserID:    feedFollowDB.UserID.String(),
		CreatedAt: feedFollowDB.CreatedAt,
		UpdatedAt: feedFollowDB.UpdatedAt,
	}
	return x, nil
}

func (apiConfig *apiConfigDefn) createFeedFollow(res http.ResponseWriter, req *http.Request) {
	type feedResDefn struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	feedRes := feedResDefn{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&feedRes)
	if err != nil {
		respondWithError(res, 200, "error decoding request body")
		return
	}
	feedFollowDB, err := apiConfig.helperFuncTCreateFeedFollow(feedRes.FeedId)
	if err != nil {
		respondWithError(res, 200, fmt.Sprintf("error while adding to database: %v", err.Error()))
		return
	}
	respondWithJson(res, 200, feedFollowDB)
}

func (apiConfig *apiConfigDefn) deleteFeedFollow(res http.ResponseWriter, req *http.Request) {
	feedFollowID := req.PathValue("feedFollowID")
	if feedFollowID == "" {
		respondWithError(res, 400, "Path Parameters not specified correctly")
		return
	}
	feedFollowUUID, err := uuid.Parse(feedFollowID)

	if err != nil {
		respondWithError(res, 400, fmt.Sprintf("path value is not valid uuid: %v :%v", feedFollowID, err.Error()))
		return
	}
	deletedRow, err := apiConfig.DB.DeleteFeedFollow(context.Background(), feedFollowUUID)
	if err != nil {
		respondWithError(res, 400, fmt.Sprintf("error deleting from database: %v", err.Error()))
		return
	}
	respondWithJson(res, 200, FeedFollowresponseStrct{
		ID:        deletedRow.ID.String(),
		FeedID:    deletedRow.FeedsID.String(),
		UserID:    deletedRow.UserID.String(),
		CreatedAt: deletedRow.CreatedAt,
		UpdatedAt: deletedRow.UpdatedAt,
	})
}

func (apiConfig *apiConfigDefn) getAllFeedFollows(res http.ResponseWriter, req *http.Request) {
	feedFollows, err := apiConfig.DB.GetFeedFollows(context.Background())
	if err != nil {
		respondWithError(res, 400, fmt.Sprintf("error retreiving from database: %v", err.Error()))
		return
	}
	returnVal := make([]FeedFollowresponseStrct, 0, len(feedFollows))
	for _, v := range feedFollows {
		returnVal = append(returnVal, FeedFollowresponseStrct{
			ID:        v.ID.String(),
			FeedID:    v.FeedsID.String(),
			UserID:    v.UserID.String(),
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}
	respondWithJson(res, 200, returnVal)
}
