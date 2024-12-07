package api

type MessageType uint

const (
	Buy MessageType = iota
	Deposit
	Delete
)

type MoneyOperationKafkaMessage struct {
	Type    MessageType `json:"type"`
	UserID  int         `json:"userID"`
	Diff    uint        `json:"diff,omitempty"`
	OrderID int         `json:"albumID,omitempty"`
}

type NotificationKafkaMessage struct {
	Type      MessageType `json:"type"`
	UserID    int         `json:"userID"`
	AlbumName string      `json:"albumName,omitempty"`
	OrderID   int         `json:"orderID,omitempty"`
	Success   *bool       `json:"success,omitempty"`
}
