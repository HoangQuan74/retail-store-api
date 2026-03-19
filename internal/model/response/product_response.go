package response

import "time"

type ProductResponse struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description,omitempty"`
	Price        float64   `json:"price"`
	Quantity     int32     `json:"quantity"`
	CategoryID   int64     `json:"category_id,omitempty"`
	CategoryName string    `json:"category_name,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ListProductResponse struct {
	Items []ProductResponse `json:"items"`
	Total int64             `json:"total"`
}
