package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoangquan/retail-store-api/pkg/model/request"
	pkgResponse "github.com/hoangquan/retail-store-api/pkg/response"
	"github.com/hoangquan/retail-store-api/services/admin/internal/service"
)

type SearchHandler struct {
	service *service.SearchService
}

func NewSearchHandler(svc *service.SearchService) *SearchHandler {
	return &SearchHandler{service: svc}
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
