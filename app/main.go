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
	"github.com/sbani/go-humanizer/numbers"
	"github.com/unrolled/secure"
)

func main() {
	loadEnvVars()

	secureMiddleware := secure.New(secure.Options{
		AllowedHostsAreRegex:  true,
		HostsProxyHeaders:     []string{"X-Forwarded-Host"},
		SSLRedirect:           true,
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:            31536000,
		STSIncludeSubdomains:  true,
		STSPreload:            true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "script-src $NONCE",
		IsDevelopment:         true,
	})

	// Create a new Gorilla mux router
	r := mux.NewRouter()
	r.Use(secureMiddleware.Handler)
	r.Use(GzipMiddleware)

	r.HandleFunc("/most_followed_users", handleMostFollowedUsers).Methods("GET")
	r.HandleFunc("/rank/{country}/{username}", handleRank).Methods("GET")

	// Enable CORS for the route
	corsMiddleware := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port number
	}

	fmt.Printf("Listening on :%s...\n", port)
	if err := http.ListenAndServe(":"+port, corsMiddleware(r)); err != nil {
		panic(err)
	}

}

/*
 * Load environment variables from the .env file and check for required env vars.
 */
func loadEnvVars() {
	// Load variables from the .env file, don't throw an error if the file doesn't exist
	godotenv.Load("./.env")

	// Check for required environment variables
	requiredEnvVars := []string{"GITHUB_API_TOKEN"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("Missing required environment variable: %s", envVar)
		}
	}
}

func handleRank(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	country := params["country"]
	username := params["username"]

	var rank int

	// Serve the response from the cache if found
	cacheKey := country + "/" + username
	userRankFinder := func() interface{} { return FindUserRank(username, country) }
	value, isCacheHit := GetFromCacheOrEvaluateFunction(cacheKey, userRankFinder)
	rank = value.(int)
	if isCacheHit {
		w.Header().Set("Served-From", "Cache")
	}

	var ordinalRank string

	if rank == -1 {
		ordinalRank = "not found"
	} else {
		ordinalRank = numbers.Ordinalize(rank)
	}

	jsonResponse, _ := json.Marshal(map[string]string{"rank": ordinalRank})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func handleMostFollowedUsers(w http.ResponseWriter, r *http.Request) {
	// Get the country from the query string after trimming any leading or trailing spaces
	country := strings.TrimSpace(r.URL.Query().Get("country"))

	// Validate the country
	isValidCountry := regexp.MustCompile(`^[a-zA-Z\s]{2,30}$`).MatchString(country)
	if !isValidCountry {
		returnJSONError(w, http.StatusBadRequest, "Invalid country name")
		return
	}

	var response []struct{ User }

	// Serve the response from the cache if found
	cacheKey := country
	mostFollowedUsersFinder := func() interface{} { return FindMostFollowedUsers(country) }
	value, isCacheHit := GetFromCacheOrEvaluateFunction(cacheKey, mostFollowedUsersFinder)
	response = value.([]struct{ User })
	if isCacheHit {
		w.Header().Set("Served-From", "Cache")
	}

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
