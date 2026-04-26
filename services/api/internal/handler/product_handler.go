package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	pkgResponse "github.com/hoangquan/retail-store-api/pkg/response"
	"github.com/hoangquan/retail-store-api/services/api/internal/service"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{service: svc}
}

// @Summary	List products
// @Tags	products
// @Produce	json
// @Param	limit query int false "Limit" default(20) minimum(1) maximum(100)
// @Param	offset query int false "Offset" default(0) minimum(0)
// @Success	200 {object} response.ProductListAPIResponse
// @Failure	500 {object} response.ErrorResponse
// @Router	/products [get]
func (h *ProductHandler) List(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 32)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 32)

	products, err := h.service.List(c.Request.Context(), int32(limit), int32(offset))
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	pkgResponse.Success(c, http.StatusOK, "Success", products)
}

// @Summary	Get product by ID
// @Tags	products
// @Produce	json
// @Param	id path int true "Product ID"
// @Success	200 {object} response.ProductAPIResponse
// @Failure	400 {object} response.ErrorResponse
// @Failure	404 {object} response.ErrorResponse
// @Router	/products/{id} [get]
func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		pkgResponse.Error(c, http.StatusNotFound, "Product not found")
		return
	}

	pkgResponse.Success(c, http.StatusOK, "Success", product)
}
