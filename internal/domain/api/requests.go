package api

type OrderActionRequest struct {
	UserID  int `json:"userID"`
	AlbumID int `json:"albumID"`
}

type MoneyOperationKafkaMessage struct {
	Type    string `json:"type"`
	UserID  int    `json:"userID"`
	Diff    uint    `json:"diff,omitempty"`
	OrderID int    `json:"albumID,omitempty"`
}

type DepositRequest struct {
	Money uint `json:"money"`
}
type BuyRequest struct {
	OrderID int `json:"money"`
}
