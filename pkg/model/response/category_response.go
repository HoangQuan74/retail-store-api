package response

import "time"

type CategoryResponse struct {
	ID        int64     `json:"id" example:"1"`
	Name      string    `json:"name" example:"Electronics"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryAPIResponse struct {
	Status  int              `json:"status" example:"200"`
	Message string           `json:"message" example:"Success"`
	Data    CategoryResponse `json:"data"`
}

type CategoryListAPIResponse struct {
	Status  int                `json:"status" example:"200"`
	Message string             `json:"message" example:"Success"`
	Data    []CategoryResponse `json:"data"`
}
