package model

import "time"

type Order struct {
	Customer   Customer     `json:"customer"`
	Date       time.Time    `json:"date"`
	TotalPrice float64      `json:"totalPrice"`
	Items      []OrderItems `json:"orderItems"`
}

type OrderItems struct {
	Album    Album   `json:"album"`
	Quantity uint    `json:"quantity"`
	Price    float64 `json:"price"`
}
