package model

import "time"

type User struct {
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

type Customer struct {
	User      User   `json:"user"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Artist struct {
	User              User   `json:"user"`
	Name              string `json:"name"`
	ProfilePictureUrl string `json:"profilePictureUrl"`
	Bio               string `json:"bio"`
}
