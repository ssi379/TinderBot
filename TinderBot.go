package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strings"
	"flag"
	"time"
)

const token string = "ENTER YOUR TOKEN HERE!"

var (
	amountOfMatches int
	verbose bool
	amountOfPeopleLiked int
)

type Person struct {
	Match bool `json:"match"`
}

type Data struct {
	Status int `json:"status"`
	Results []Prospect `json:"results"`
}

type Prospect struct {
	Id string `json:"_id"`
	Name string `json:"name"`
	Birthdate string `json:"birth_date"`
}

func main() {
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.Parse()

	fmt.Println("Running...")

	// Print stats every 5 minutes
	go statsScheduler()

	// Get prospects
	for {
		getProspects()
	}
}

// Sends requests to tinder for Prospects and likes all of them
func getProspects() {
	const url string = "https://api.gotinder.com/user/recs"

	// Create HTTP Client
	httpClient := http.Client{}

	// Create POST request
	req, err := http.NewRequest("POST", url, strings.NewReader("{\"limit\":40}"))

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Set Headers
	setHeaders(req)

	// Execute request
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Handle response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Sometimes Tinder returns an empty array
	if len(body) < 1 {
		return
	}

	// Parse JSON
	var prospects Data
	err = json.Unmarshal(body, &prospects)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//Like all the new prospects
	results := prospects.Results
	handleProspects(results)
}

func handleProspects(prospects []Prospect) {
	for _, prospect := range prospects {
		go likeUser(prospect.Id)
		if verbose {
			fmt.Printf("Liking %s\n", prospect.Name)
		}
	}
}

func likeUser(userId string) {
	const baseURL string = "https://api.gotinder.com/like/"

	url := baseURL + userId

	// Create HTTP Client
	httpClient := http.Client{}

	// Create GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Set Headers
	setHeaders(req)

	// Execute request
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	amountOfPeopleLiked++

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Check response for match
	var p Person
	err = json.Unmarshal(body, &p)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("At line 176")
	}

	// TODO: Test and fix this
	if p.Match {
		// We have a match! Increment counter
		amountOfMatches++
		if (verbose) {
			fmt.Println("We have a new match!")
		}
	}
}

// Sets the request headers
func setHeaders(req *http.Request) (*http.Request) {
	// Set Headers
	req.Header.Set("platform","android")
	req.Header.Set("User-Agent","Tinder Android Version 3.2.1")
	req.Header.Set("X-Auth-Token", token)
	req.Header.Set("os-version","19")
	req.Header.Set("app-version","759")
	return req
}

func statsScheduler() {
	for _ = range time.Tick(5 * time.Minute) {
		printStats()
	}
}

// Prints the stats
func printStats() {
	fmt.Printf("We have liked %v prospects, we gained %v new matches.", amountOfPeopleLiked, amountOfMatches)
}
