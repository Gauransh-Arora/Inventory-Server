package models

import "time"

type Product struct{
	ID int `json:"id"`
	Barcode string `json:"barcode"`
	ItemName string `json:"item_name"`
	MRP float64 `json:"mrp"`
	CreatedAt time.Time `json:"created_at"`
}