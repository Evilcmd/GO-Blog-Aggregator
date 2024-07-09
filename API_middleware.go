package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
)

func (apiConfig *apiConfigDefn) authenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		apiKey := req.Header.Get("Authorization")
		if len(apiKey) == 0 {
			respondWithError(res, 400, "authorization header not set")
			return
		}
		apiKey = strings.Split(apiKey, " ")[1]

		dbUser, err := apiConfig.DB.GetUserByAPIKey(context.Background(), sql.NullString{
			String: apiKey,
			Valid:  true,
		})
		if err != nil {
			respondWithError(res, 400, fmt.Sprintf("error while fetching user: %v", err.Error()))
			return
		}

		apiConfig.user = dbUser

		next.ServeHTTP(res, req)
	}
}
