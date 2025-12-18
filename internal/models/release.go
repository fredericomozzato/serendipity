package models

type Release struct {
	ID        int      `json:"id"`
	Title     string   `json:"title"`
	Artists   []Artist `json:"artists"`
	Year      int      `json:"year"`
	Genres    []string `json:"genres"`
	Country   string   `json:"country"`
	Images    []Image  `json:"images,omitempty"`
	Tracklist []Track  `json:"tracklist"`
	URI       string   `json:"uri"`
	Videos    []Video  `json:"videos,omitempty"`
}
