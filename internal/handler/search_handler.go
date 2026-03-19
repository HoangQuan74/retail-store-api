package handler

import (
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/kainguyen/retail-store-api/internal/model/request"
	"github.com/kainguyen/retail-store-api/internal/service"
	pkgResponse "github.com/kainguyen/retail-store-api/pkg/response"
)

type SearchHandler struct {
	service *service.SearchService
}

func NewSearchHandler(client *elasticsearch.Client, indexName string) *SearchHandler {
	return &SearchHandler{service: service.NewSearchService(client, indexName)}
}

func (h *SearchHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/api/v1/search/products", h.SearchProducts)
}

func (h *SearchHandler) SearchProducts(c *gin.Context) {
	var req request.SearchProductRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.SearchProducts(req.Q, req.Limit, req.Offset)
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	pkgResponse.Success(c, http.StatusOK, "Success", gin.H{
		"items": result.Items,
		"total": result.Total,
	})
}
