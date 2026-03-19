package response

import "time"

type ProductResponse struct {
	ID           int64     `json:"id" example:"1"`
	Name         string    `json:"name" example:"Laptop"`
	Description  string    `json:"description,omitempty" example:"A powerful laptop"`
	Price        float64   `json:"price" example:"999.99"`
	Quantity     int32     `json:"quantity" example:"10"`
	CategoryID   int64     `json:"category_id,omitempty" example:"1"`
	CategoryName string    `json:"category_name,omitempty" example:"Electronics"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ProductAPIResponse struct {
	Status  int             `json:"status" example:"200"`
	Message string          `json:"message" example:"Success"`
	Data    ProductResponse `json:"data"`
}

type ProductListAPIResponse struct {
	Status  int               `json:"status" example:"200"`
	Message string            `json:"message" example:"Success"`
	Data    []ProductResponse `json:"data"`
}
