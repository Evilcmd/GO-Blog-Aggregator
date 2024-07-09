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

type DBUser struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Apikey    string    `json:"api_key"`
}

func createUserFromDBUser(usr database.User) DBUser {
	dbusr := DBUser{
		usr.ID,
		usr.CreatedAt,
		usr.UpdatedAt,
		usr.Name,
		usr.Apikey.String,
	}
	return dbusr
}

func (apiConfig *apiConfigDefn) createUser(res http.ResponseWriter, req *http.Request) {
	type userReqDefn struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(req.Body)
	userReq := userReqDefn{}
	err := decoder.Decode(&userReq)
	if err != nil {
		respondWithError(res, 400, "not able to decode user creation request")
		return
	}

	dbUser, err := apiConfig.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      userReq.Name,
	})
	if err != nil {
		respondWithError(res, 400, fmt.Sprintf("error creating user: %v", err.Error()))
	}
	respondWithJson(res, 200, createUserFromDBUser(dbUser))
}

func (apiConfig *apiConfigDefn) getUserByApiKey(res http.ResponseWriter, req *http.Request) {
	respondWithJson(res, 200, createUserFromDBUser(apiConfig.user))
}
