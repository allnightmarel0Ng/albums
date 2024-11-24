package model

import "time"

type Album struct {
	Artist      Artist    `json:"artist"`
	Name        string    `json:"name"`
	ReleaseDate time.Time `json:"releaseDate"`
	CoverArtUrl string    `json:"coverArtUrl"`
	Price       float64   `json:"price"`
	Genre       string    `json:"genre"`
	Tracks      []Track   `json:"tracks"`
}

type Track struct {
	Name         string `json:"name"`
	Duration     uint   `json:"duration"`
	AudioFileUrl string `json:"audioFileUrl"`
}
