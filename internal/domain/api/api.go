package api

import (
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/domain/model"
)

type Response interface {
	GetCode() int
	GetKeyValue() (string, interface{})
}

type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func (e *ErrorResponse) GetCode() int {
	return e.Code
}

func (e *ErrorResponse) GetKeyValue() (string, interface{}) {
	return "error", e.Error
}

type AuthorizationResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error,omitempty"`
	Jwt   string `json:"jwt,omitempty"`
}

func (a *AuthorizationResponse) GetCode() int {
	return a.Code
}

func (a *AuthorizationResponse) GetKeyValue() (string, interface{}) {
	switch a.Code {
	case http.StatusOK:
		return "jwt", a.Jwt
	default:
		return "error", a.Error
	}
}

type AlbumsResponse struct {
	Code  int           `json:"code"`
	Error string        `json:"error,omitempty"`
	Data  []model.Album `json:"data,omitempty"`
}

func (a *AlbumsResponse) GetCode() int {
	return a.Code
}

func (a *AlbumsResponse) GetKeyValue() (string, interface{}) {
	switch a.Code {
	case http.StatusOK:
		return "jwt", a.Data
	default:
		return "error", a.Error
	}
}

type UserProfileResponse struct {
	Code  int        `json:"code"`
	Error string     `json:"error,omitempty"`
	User  model.User `json:"user,omitempty"`
}

func (u *UserProfileResponse) GetCode() int {
	return u.Code
}

func (a *UserProfileResponse) GetKeyValue() (string, interface{}) {
	switch a.Code {
	case http.StatusOK:
		return "user", a.User
	default:
		return "error", a.Error
	}
}

type ArtistProfileResponse struct {
	Code   int           `json:"code"`
	Error  string        `json:"error,omitempty"`
	Albums []model.Album `json:"albums,omitempty"`
}

func (a *ArtistProfileResponse) GetCode() int {
	return a.Code
}

func (a *ArtistProfileResponse) GetKeyValue() (string, interface{}) {
	switch a.Code {
	case http.StatusOK:
		return "albums", a.Albums
	default:
		return "error", a.Error
	}
}

type JWTClaims struct {
	ID      int
	IsAdmin bool
	Exp     int64
}
