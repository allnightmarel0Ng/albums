package model

import "time"

type User struct {
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordHash"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
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
