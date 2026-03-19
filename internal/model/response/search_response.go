package response

import "github.com/kainguyen/retail-store-api/pkg/elasticsearch"

type SearchProductResponse struct {
	Items []elasticsearch.ProductDocument `json:"items"`
	Total int64                           `json:"total" example:"42"`
}

type SearchProductAPIResponse struct {
	Status  int                   `json:"status" example:"200"`
	Message string                `json:"message" example:"Success"`
	Data    SearchProductResponse `json:"data"`
}
