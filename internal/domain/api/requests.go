package api

type OrderActionRequest struct {
	UserID  int `json:"userID" binding:"required"`
	AlbumID int `json:"albumID" binding:"required"`
}

type DepositRequest struct {
	Money uint `json:"money" binding:"required"`
}

type SearchRequest struct {
	Query string `json:"query" binding:"required"`
}

type RandomEntitiesRequest struct {
	ArtistsCount uint `json:"artistsCount" binding:"required"`
	AlbumsCount  uint `json:"albumsCount" binding:"required"`
}

type RegistrationRequest struct {
	Email    string `json:"email" binding:"required"`
	IsAdmin  *bool  `json:"isAdmin" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	ImageURL string `json:"imageURL" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type NotificationSubscribeRequest struct {
	Jwt string `json:"jwt" binding:"required"`
}
