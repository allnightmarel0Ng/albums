package api

type OrderActionRequest struct {
	UserID  int `json:"userID"`
	AlbumID int `json:"albumID"`
}

type MoneyOperationKafkaMessage struct {
	Type    string `json:"type"`
	UserID  int    `json:"userID"`
	Diff    int    `json:"diff,omitempty"`
	AlbumID int    `json:"albumID,omitempty"`
}
