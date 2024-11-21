package model

import "time"

type Album struct {
	Artist      Artist    `json:"artist"`
	Name        string    `json:"name"`
	ReleaseDate time.Time `json:"releaseDate"`
	CoverArtUrl string    `json:"coverArtUrl"`
	Price       float64   `json:"price"`
	Genre       string    `json:"genre"`
}

type Track struct {
	Album        Album  `json:"album"`
	Name         string `json:"name"`
	Number       uint   `json:"number"`
	Duration     uint   `json:"duration"`
	AudioFileUrl string `json:"audioFileUrl"`
}
