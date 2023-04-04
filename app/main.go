package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	loadEnvVars()

	// Create a new Gorilla mux router
	r := mux.NewRouter()

	r.HandleFunc("/most_followed_users", handleMostFollowedUsers).Methods("GET")

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

func handleMostFollowedUsers(w http.ResponseWriter, r *http.Request) {
	// Get the country from the query string after trimming any leading or trailing spaces
	country := strings.TrimSpace(r.URL.Query().Get("country"))

	// Validate the country
	isValidCountry := regexp.MustCompile(`^[a-zA-Z\\s]{2,30}$`).MatchString(country)
	if !isValidCountry {
		returnJSONError(w, http.StatusBadRequest, "Invalid country name")
		return
	}

	// Call the GitHub API
	var response []struct{ User } = FindMostFollowedUsers(country)

	// Convert the response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		returnJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Set the content type header to indicate that the response body is JSON
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func returnJSONError(w http.ResponseWriter, statusCode int, message string) {

	type ErrorResponse struct {
		Message string `json:"message"`
	}

	// Return an error response in JSON format
	jsonError, err := json.Marshal(ErrorResponse{Message: message})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonError)
}
