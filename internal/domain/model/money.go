package model

import "time"

type Order struct {
	ID         int       `json:"id"`
	Orderer    User      `json:"orderer"`
	Date       time.Time `json:"date"`
	TotalPrice float64   `json:"totalPrice"`
	IsPaid     bool      `json:"isPaid"`
	Albums     []Album   `json:"albums"`
}
