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
