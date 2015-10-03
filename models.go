package main

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

type FacebookToken struct {
	Token string `json:"facebook_token"`
}

type PingWrapper struct {
	Lon float32 `json:"lon"`
	Lat float32 `json:"lat"`
}

type Configuration struct {
	FacebookToken string `json:"facebook_token"`
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

type Profile struct {
	TinderAccessToken string `json:"token"`
	User User `json:"user"`
}

type User struct {
	Id string `json:"_id"`
	Name string `json:"full_name"`
	TinderAccessToken string `json:"api_token"`
}
