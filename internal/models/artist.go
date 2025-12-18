package models

type Artist struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"resource_url"`
}
