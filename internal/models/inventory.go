package models

import "time"

type InventoryLog struct {
	ID        int       `json:"id"`
	ProductID int       `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Updated   bool      `json:"updated"`
	CreatedAt time.Time `json:"created_at"`
}


type LogFilter struct {
	Updated   *bool
	ProductID *int
	DateFrom  *time.Time
	DateTo    *time.Time
}


type PaginatedLogs struct {
	Data       []InventoryLog `json:"data"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalCount int            `json:"total_count"`
	TotalPages int            `json:"total_pages"`
}
