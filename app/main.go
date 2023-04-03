package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	loadEnvVars()

	// Create a new Gorilla mux router
	r := mux.NewRouter()

	// Define the route for the HTTP handler function
	r.HandleFunc("/most_followed_users", func(w http.ResponseWriter, r *http.Request) {
		// Get the country from the query string
		country := r.URL.Query().Get("country")

		// Find the top most followed users in the given country
		githubResp := FindMostFollowedUsers(country)

		// Convert the response to a JSON string
		jsonBytes, err := json.Marshal(githubResp.Data.Search.Nodes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the response to the HTTP response writer
		w.Header().Add("Content-Type", "application/json")
		w.Write(jsonBytes)
	})

	// Enable CORS for the route
	corsMiddleware := handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}),
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port number
	}

	fmt.Printf("Listening on :%s...\n", port)
	if err := http.ListenAndServe(":"+port, corsMiddleware(r)); err != nil {
		panic(err)
	}

}

func loadEnvVars() {
	// Load from the .env file
	err := godotenv.Load("./.env")
	isProduction := os.Getenv("PRODUCTION")
	if err != nil && isProduction != "true" {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
