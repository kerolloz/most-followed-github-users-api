package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type User struct {
	Name      string `json:"name"`
	Login     string `json:"login"`
	Bio       string `json:"bio"`
	Followers struct {
		TotalCount int `json:"totalCount"`
	} `json:"followers"`
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

func FindMostFollowedUsers(country string) GithubResponse {
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

	GitHubAPIToken := os.Getenv("GITHUB_API_TOKEN")
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
