package models

import "time"

type InventoryLog struct{
	ID int `json:"id"`
	ProductID int `json:"product_id"`
	Quantity int `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}