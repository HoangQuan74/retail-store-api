package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kainguyen/retail-store-api/internal/app"
	"github.com/kainguyen/retail-store-api/internal/model/request"
	"github.com/kainguyen/retail-store-api/internal/service"
	pkgResponse "github.com/kainguyen/retail-store-api/pkg/response"
)

type SearchHandler struct {
	service *service.SearchService
}

func NewSearchHandler(ctx *app.AppContext, router *gin.Engine) *SearchHandler {
	h := &SearchHandler{service: service.NewSearchService(ctx.ESClient, ctx.ProductIndex)}
	router.GET("/api/v1/search/products", h.SearchProducts)
	return h
}

// @Summary	Search products
// @Tags	search
// @Produce	json
// @Param	q query string true "Search query" minLength(1)
// @Param	limit query int false "Limit" default(20) minimum(1) maximum(100)
// @Param	offset query int false "Offset" default(0) minimum(0)
// @Success	200 {object} response.SearchProductAPIResponse
// @Failure	400 {object} response.ErrorResponse
// @Failure	500 {object} response.ErrorResponse
// @Router	/search/products [get]
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
