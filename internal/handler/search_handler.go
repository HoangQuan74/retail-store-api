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
