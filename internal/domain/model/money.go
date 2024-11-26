package model

import "time"

type Order struct {
	ID         int         `json:"id"`
	Orderer    User        `json:"orderer"`
	Date       time.Time   `json:"date"`
	TotalPrice float64     `json:"totalPrice"`
	Items      []OrderItem `json:"orderItems"`
}

type OrderItem struct {
	ID    int   `json:"id"`
	Album Album `json:"album"`
}
