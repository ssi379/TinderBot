package main

import (
	"fmt"
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/json"
)

var (
	token string
)

func main() {
	for len(token) < 1 {
		fmt.Println("Please enter your Tinder Access Token: ")
		fmt.Scanln(&token)
	}

	for {
		getProspects()
	}
}

// Sends requests to Tinder for new prospects and likes all of them.
func getProspects() {
	const url string = "https://api.gotinder.com/user/recs"

	// Create HTTP client
	httpClient := http.Client{}

	// Create POST request
	req, err := http.NewRequest("POST", url, strings.NewReader("{\"limit\":40}"))
	if err != nil {
		fmt.Println("Error creating POST request: " + err.Error())
		return
	}

	// Set Headers
	setHeaders(req)

	// Execute the request
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error retrieving prospects: " + err.Error())
		return
	}

	// Handle response
	body, err := ioutil.ReadAll(resp.Body);
	if err != nil {
		fmt.Println("Error reading response from Tinder: " + err.Error())
		return
	}

	// Sometimes Tinder returns an empty array for some reason.
	if len(body) < 1 {
		fmt.Println("Tinder returned no results.")
		return;
	}

	// Unmarshal JSON
	var prospects Data
	if err := json.Unmarshal(body, &prospects); err != nil {
		fmt.Println("Error unmarshaling JSON: " + err.Error())
		return
	}

	// Like all new prospects
	likeAllProspects(prospects.Results)
}

func likeAllProspects(prospects []Prospect) {
	for _, prospect := range prospects {
		go likeUser(prospect)
	}
}

func likeUser(prospect Prospect) {
	const baseURL string = "https://api.gotinder.com/like/"

	// Build url
	url := baseURL + prospect.Id

	// Create HTTP Client
	httpClient := http.Client{}

	// Create GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating GET request: " + err.Error())
		return
	}

	// Set Headers
	setHeaders(req)

	// Execute request
	if _, err := httpClient.Do(req); err != nil {
		fmt.Println("Error executing GET request: " + err.Error())
		return
	}

	fmt.Println("Succesfully liked: " + prospect.Name + " Birthdate: " + prospect.Birthdate + " UserId: " + prospect.Id)
	// TODO: Check if match
}

// Sets the request headers
func setHeaders(req *http.Request) (*http.Request) {
	// Set Headers
	req.Header.Set("platform", "android")
	req.Header.Set("User-Agent", "Tinder Android Version 3.2.1")
	req.Header.Set("X-Auth-Token", token)
	req.Header.Set("os-version", "19")
	req.Header.Set("app-version", "759")
	return req
}
