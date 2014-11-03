package main

import (
	"fmt"
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/json"
	"errors"
	"log"
	"bytes"
	"os"
	"bufio"
)

var (
	config Configuration
	amountOfEmptyResults int = 0
	tinderAccessToken string
)

func main() {
	fmt.Println("Welcome to TinderBot")

	if !checkConfig() {
		for len(config.FacebookToken) < 1 {
			fmt.Println("Please enter your Facebook Access Token: ")
			fmt.Scanln(&config.FacebookToken)
		}
		config.save()
	} else {
		if err := loadConfig(); err != nil {
			log.Fatal("TinderBot is unable to load the configuration file")
		}
	}

	if err := retrieveAccessToken(); err != nil {
		log.Fatal("Failed to retrieve Access Token.\nIs the token you supplied valid?")
	}

	inputReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Do you want to change your location? [Y/n]")
		input, err := inputReader.ReadString('\n')

		if err != nil {
			fmt.Println("There were errors reading your input, exiting program.")
			return
		}

		if strings.ContainsAny(input, "Y y") {
			spoofLocation()
		} else {
			amountOfEmptyResults = 0
		}

		for amountOfEmptyResults < 10 {
			getProspects()
		}

		fmt.Println("We haven't retrieved any results in a while...")
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
		amountOfEmptyResults++;
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

	// Reset amount of empty results since we still are able to like people
	if (len(prospects.Results) > 0) {
		amountOfEmptyResults = 0;
	}
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
	req.Header.Set("User-Agent", "Tinder Android Version 3.3.2")
	req.Header.Set("X-Auth-Token", tinderAccessToken)
	req.Header.Set("os-version", "19")
	req.Header.Set("app-version", "763")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	return req
}

func retrieveAccessToken() error {
	if len(config.FacebookToken) > 0 {

		facebookToken := FacebookToken{Token: config.FacebookToken}

		// Marshal struct into JSON
		b, err := json.Marshal(facebookToken)
		if err != nil {
			return err;
		}

		// Create IO reader
		postBody := bytes.NewReader(b)

		// Make request
		const authURL = "https://api.gotinder.com/auth"
		resp, err := http.Post(authURL, "application/json", postBody)

		response, err := ioutil.ReadAll(resp.Body)

		// Unmarshal response
		var profile Profile
		if err := json.Unmarshal(response, &profile); err != nil {
			fmt.Println("Error unmarshaling Tinder's response")
			return err
		}

		tinderAccessToken = profile.TinderAccessToken

		fmt.Println("Authenticated " + profile.User.Name + " with Tinder")

		return nil
	} else {
		return errors.New("Invalid Facebook token")
	}
}

func checkConfig() bool {
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func loadConfig() error {
	b, _ := ioutil.ReadFile("config.json")
	if err := json.Unmarshal(b, &config); err != nil {
		return err
	}
	return nil
}

func (c *Configuration) save() error {
	output, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		fmt.Println("Unable to marshal configuration")
	}
	return ioutil.WriteFile("config.json", output, 0600)
}

func spoofLocation() {
	fmt.Println("Please enter your new Latitude: ")
	fmt.Scanf("%f\n", &config.Lat)
	fmt.Println("Please enter your new Longitude: ")
	fmt.Scanf("%f\n", &config.Lon)
	fmt.Printf("Setting location to: %f, %f\n", config.Lat, config.Lon)

	// Marshal new location into json
	ping := PingWrapper{Lat: config.Lat, Lon: config.Lon}
	b, err := json.Marshal(ping)
	if err != nil {
		fmt.Println("Error marshaling JSON: " + err.Error())
		return
	}

	// Create IO reader
	postBody := bytes.NewReader(b)

	// Create HTTP Client
	httpClient := http.Client{}

	// Make request
	const pingURL = "https://api.gotinder.com/user/ping"
	req, err := http.NewRequest("POST", pingURL, postBody)
	if err != nil {
		fmt.Println("Error creating POST request: " + err.Error())
		return
	}

	// Set headers
	setHeaders(req)

	// Submit Ping request
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Setting new location failed")
	}

	if resp.StatusCode == 200 {
		fmt.Println("New location succesfully set.")
		amountOfEmptyResults = 0;
	}
}
