package api

import (
	"github.com/allnightmarel0Ng/albums/internal/domain/model"
)

type Response interface {
	GetCode() int
}

type ErrorResponse struct {
	Code  int    `json:"-"`
	Error string `json:"error"`
}

func (e *ErrorResponse) GetCode() int {
	return e.Code
}

type AuthenticationResponse struct {
	Code  int    `json:"-"`
	Error string `json:"error,omitempty"`
	Jwt   string `json:"jwt,omitempty"`
}

func (a *AuthenticationResponse) GetCode() int {
	return a.Code
}

type AlbumsResponse struct {
	Code  int           `json:"-"`
	Error string        `json:"error,omitempty"`
	Data  []model.Album `json:"data,omitempty"`
}

func (a *AlbumsResponse) GetCode() int {
	return a.Code
}

type UserProfileResponse struct {
	Code  int        `json:"-"`
	Error string     `json:"error,omitempty"`
	User  model.User `json:"user,omitempty"`
}

func (u *UserProfileResponse) GetCode() int {
	return u.Code
}

type ArtistProfileResponse struct {
	Code   int           `json:"-"`
	Error  string        `json:"error,omitempty"`
	Albums []model.Album `json:"albums,omitempty"`
}

func (a *ArtistProfileResponse) GetCode() int {
	return a.Code
}

type AuthorizationResponse struct {
	Code    int    `json:"-"`
	Error   string `json:"error,omitempty"`
	ID      int    `json:"id,omitempty"`
	IsAdmin bool   `json:"isAdmin,omitempty"`
}

func (j *AuthorizationResponse) GetCode() int {
	return j.Code
}

type OrderActionResponse struct {
	Code  int    `json:"-"`
	Error string `json:"error,omitempty"`
}

func (o *OrderActionResponse) GetCode() int {
	return o.Code
}

type UserOrdersResponse struct {
	Code   int           `json:"-"`
	Error  string        `json:"error,omitempty"`
	Orders []model.Order `json:"orders,omitempty"`
}

func (u *UserOrdersResponse) GetCode() int {
	return u.Code
}
