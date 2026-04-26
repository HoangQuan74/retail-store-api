package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	pkgResponse "github.com/hoangquan/retail-store-api/pkg/response"
	"github.com/hoangquan/retail-store-api/services/api/internal/service"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: svc}
}

// @Summary	List all categories
// @Tags	categories
// @Produce	json
// @Success	200 {object} response.CategoryListAPIResponse
// @Failure	500 {object} response.ErrorResponse
// @Router	/categories [get]
func (h *CategoryHandler) List(c *gin.Context) {
	categories, err := h.service.List(c.Request.Context())
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	pkgResponse.Success(c, http.StatusOK, "Success", categories)
}

// @Summary	Get category by ID
// @Tags	categories
// @Produce	json
// @Param	id path int true "Category ID"
// @Success	200 {object} response.CategoryAPIResponse
// @Failure	400 {object} response.ErrorResponse
// @Failure	404 {object} response.ErrorResponse
// @Router	/categories/{id} [get]
func (h *CategoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	category, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		pkgResponse.Error(c, http.StatusNotFound, "Category not found")
		return
	}

	pkgResponse.Success(c, http.StatusOK, "Success", category)
}
