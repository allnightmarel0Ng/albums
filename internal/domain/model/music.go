package model

type Artist struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Genre    string `json:"genre"`
	ImageURL string `json:"imageURL"`
}

type Album struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Author   Artist  `json:"author,omitempty"`
	ImageURL string  `json:"imageURL"`
	Price    float64 `json:"price"`
	Tracks   []Track `json:"tracks"`
}

type Track struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Number int    `json:"number"`
}
