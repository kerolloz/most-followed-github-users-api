package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type User struct {
	Name      string `json:"name"`
	Login     string `json:"login"`
	Bio       string `json:"bio"`
	Followers struct {
		TotalCount int `json:"totalCount"`
	} `json:"followers"`
	AvatarURL string `json:"avatarUrl"`
}

type GithubResponse struct {
	Data struct {
		Search struct {
			Nodes []struct {
				User
			} `json:"nodes"`
		} `json:"search"`
	} `json:"data"`
}

func main() {
	// Load the environment variables from the .env file
	err := godotenv.Load("./.env")
	isProduction := os.Getenv("PRODUCTION")
	if err != nil && isProduction != "true" {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port number
	}

	// Create a new Gorilla mux router
	r := mux.NewRouter()

	// Define the route for the HTTP handler function
	r.HandleFunc("/most_followed_users", func(w http.ResponseWriter, r *http.Request) {
		// Get the country from the query string
		country := r.URL.Query().Get("country")

		// Find the top most followed users in the given country
		githubResp := findMostFollowedUsersInCountry(country)

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

	fmt.Printf("Listening on :%s...\n", port)
	if err := http.ListenAndServe(":"+port, corsMiddleware(r)); err != nil {
		panic(err)
	}

}

func findMostFollowedUsersInCountry(country string) GithubResponse {
	GitHubAPIToken := os.Getenv("GITHUB_API_TOKEN")
	// Set the GraphQL query to retrieve the top 10 most followed users in Egypt
	query := fmt.Sprintf(`
		{
			search(query: "location:%s sort:followers-desc", type: USER, first: 50) {
				nodes {
					... on User {
						name
						login
						bio
						followers {
							totalCount
						}
						avatarUrl
					}
				}
			}
		}
`, country)

	// Define the GraphQL request payload
	payload := map[string]string{
		"query": query,
	}

	// Convert the payload to a JSON string
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(payloadBytes))
	if err != nil {
		panic(err)
	}

	// Set the authorization header to include the access token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", GitHubAPIToken))

	// Set the content type header to indicate that the request body is JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a new HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Parse the JSON response into a GithubResponse struct
	var githubResp GithubResponse
	if err := json.Unmarshal(body, &githubResp); err != nil {
		panic(err)
	}
	return githubResp
}
