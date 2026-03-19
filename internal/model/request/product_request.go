package request

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Quantity    int32   `json:"quantity" binding:"gte=0"`
	CategoryID  int64   `json:"category_id"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Quantity    int32   `json:"quantity" binding:"gte=0"`
	CategoryID  int64   `json:"category_id"`
}

type ListProductRequest struct {
	Limit  int32 `form:"limit,default=20" binding:"gte=1,lte=100"`
	Offset int32 `form:"offset,default=0" binding:"gte=0"`
}
