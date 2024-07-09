package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Evilcmd/GO-Blog-Aggregator/internal/database"
)

func (apiConfig *apiConfigDefn) GetUserPosts(res http.ResponseWriter, req *http.Request) {
	posts, err := apiConfig.DB.GetPostsByUser(context.Background(), database.GetPostsByUserParams{
		UserID: apiConfig.user.ID,
		Limit:  2,
	})
	// fmt.Println(apiConfig.user.ID)
	if err != nil {
		respondWithError(res, 400, fmt.Sprintf("error in getting posts: %v\n", err.Error()))
		return
	}
	respondWithJson(res, 200, posts)
}
