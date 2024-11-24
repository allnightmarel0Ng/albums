package model

import (
	"net/http"
)

type Response interface {
	GetCode() int
	GetKeyValue() (string, interface{})
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
	Code  int     `json:"code"`
	Error string  `json:"error,omitempty"`
	Data  []Album `json:"data,omitempty"`
}

func (a *AlbumsResponse) GetCode() int {
	return a.Code
}

func (a *AlbumsResponse) GetKeyValue() (string, interface{}) {
	switch a.Code {
	case http.StatusOK:
		return "data", a.Data
	default:
		return "error", a.Error
	}
}
