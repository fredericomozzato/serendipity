package models

type Image struct {
	ID   int
	Type string `json:"type"`
	URI  string `json:"uri"`
}
