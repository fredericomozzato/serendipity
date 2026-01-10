package models

type Release struct {
	ID        int
	Title     string
	Artists   []Artist
	Year      int
	Genres    []string
	Country   string
	Images    []Image
	Tracklist []Track
	URI       string
	Videos    []Video
}
