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

type UserProfileResponse struct {
	Code      int           `json:"-"`
	User      model.User    `json:"user"`
	Purchased []model.Album `json:"purchasedAlbums,omitempty"`
}

func (u *UserProfileResponse) GetCode() int {
	return u.Code
}

type ArtistProfileResponse struct {
	Code   int           `json:"-"`
	Artist model.Artist  `json:"artist"`
	Albums []model.Album `json:"albums,omitempty"`
}

func (a *ArtistProfileResponse) GetCode() int {
	return a.Code
}

type AlbumProfileResponse struct {
	Code  int         `json:"-"`
	Album model.Album `json:"album"`
}

func (a *AlbumProfileResponse) GetCode() int {
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

type UserOrdersResponse struct {
	Code   int           `json:"-"`
	Orders []model.Order `json:"orders,omitempty"`
}

func (u *UserOrdersResponse) GetCode() int {
	return u.Code
}

type UnpaidUserOrderResponse struct {
	Code  int         `json:"-"`
	Error string      `json:"error,omitempty"`
	Order model.Order `json:"order,omitempty"`
}

func (u *UnpaidUserOrderResponse) GetCode() int {
	return u.Code
}

type SearchEngineResponse struct {
	Code    int            `json:"-"`
	Artists []model.Artist `json:"artists,omitempty"`
	Albums  []model.Album  `json:"albums,omitempty"`
}

func (s *SearchEngineResponse) GetCode() int {
	return s.Code
}

type BuyLogsResponse struct {
	Code      int            `json:"-"`
	Logs      []model.BuyLog `json:"logs"`
	LogsCount uint           `json:"logsCount"`
}

func (b *BuyLogsResponse) GetCode() int {
	return b.Code
}
