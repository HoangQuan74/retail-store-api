package response

import "github.com/kainguyen/retail-store-api/pkg/elasticsearch"

type SearchProductResponse struct {
	Items []elasticsearch.ProductDocument `json:"items"`
	Total int64                           `json:"total"`
}
