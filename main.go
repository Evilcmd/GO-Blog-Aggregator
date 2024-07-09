package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Evilcmd/GO-Blog-Aggregator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfigDefn struct {
	DB   *database.Queries
	user database.User
}

func main() {
	godotenv.Load()
	PORT := os.Getenv("PORT")

	dbURL := os.Getenv("DBURL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)
	apiConfig := apiConfigDefn{dbQueries, database.User{}}

	mux := http.NewServeMux()

	// Server Check
	mux.HandleFunc("GET /v1/healthz", healthz)
	mux.HandleFunc("GET /v1/err", erreHandle)

	// create user
	mux.HandleFunc("POST /v1/users", apiConfig.createUser)

	// Get user by using API key
	mux.HandleFunc("GET /v1/users", apiConfig.authenticationMiddleware(apiConfig.getUserByApiKey))

	// handler to create a feed
	mux.HandleFunc("POST /v1/feeds", apiConfig.authenticationMiddleware(apiConfig.createFeeds))

	// get all the feeds
	mux.HandleFunc("GET /v1/feeds", apiConfig.GetAllFeeds)

	// add a feed follow
	mux.HandleFunc("POST /v1/feed_follows", apiConfig.authenticationMiddleware(apiConfig.createFeedFollow))

	// delete a feed follow
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", apiConfig.authenticationMiddleware(apiConfig.deleteFeedFollow))

	// get all feed follows
	mux.HandleFunc("GET /v1/feed_follows", apiConfig.authenticationMiddleware(apiConfig.getAllFeedFollows))

	server := http.Server{
		Addr:    fmt.Sprintf(":%v", PORT),
		Handler: mux,
	}
	fmt.Printf("Starting Server on :%v\n", PORT)
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server")
	}
}
