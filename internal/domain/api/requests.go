package api

type OrderActionRequest struct {
	UserID  int `json:"userID"`
	AlbumID int `json:"albumID"`
}
