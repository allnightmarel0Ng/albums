package model

type User struct {
	ID        int     `json:"id"`
	Email     string  `json:"email"`
	IsAdmin   bool    `json:"isAdmin"`
	Nickname  string  `json:"nickname"`
	Balance   float64 `json:"balance"`
	ImageURL  string  `json:"imageURL"`
	Purchased []Album `json:"purchasedAlbums"`
}
