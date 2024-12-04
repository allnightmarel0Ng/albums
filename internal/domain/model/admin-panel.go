package model

import "time"

type BuyLog struct {
	ID          int       `json:"id"`
	Buyer       User      `json:"buyer"`
	Album       Album     `json:"album"`
	LoggingTime time.Time `json:"loggingTime"`
}
