package models

type Product struct {
	ID       int     `json:"id"`
	ICode    int     `json:"icode"`
	ItemName string  `json:"item_name"`
	BatchNo  int     `json:"batch_no"`
	MRP      float64 `json:"mrp"`
	Barcode  string  `json:"barcode"`
}

type Pagination struct {
	Page     int 
	PageSize int 
}

type PaginatedProducts struct {
	Data       []Product `json:"data"`
	Page       int       `json:"page"`
	PageSize   int       `json:"page_size"`
	TotalCount int       `json:"total_count"`
	TotalPages int       `json:"total_pages"`
}