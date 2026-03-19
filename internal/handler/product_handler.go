package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kainguyen/retail-store-api/internal/app"
	"github.com/kainguyen/retail-store-api/internal/model/request"
	"github.com/kainguyen/retail-store-api/internal/service"
	pkgResponse "github.com/kainguyen/retail-store-api/pkg/response"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(ctx *app.AppContext, router *gin.Engine) *ProductHandler {
	h := &ProductHandler{service: service.NewProductService(ctx.Queries, ctx.Publisher)}
	products := router.Group("/api/v1/products")
	{
		products.POST("", h.Create)
		products.GET("", h.List)
		products.GET("/:id", h.GetByID)
		products.PUT("/:id", h.Update)
		products.DELETE("/:id", h.Delete)
	}
	return h
}

// @Summary	Create a product
// @Tags	products
// @Accept	json
// @Produce	json
// @Param	body body request.CreateProductRequest true "Product payload"
// @Success	201 {object} response.ProductAPIResponse
// @Failure	400 {object} response.ErrorResponse
// @Failure	500 {object} response.ErrorResponse
// @Router	/products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var req request.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	pkgResponse.Success(c, http.StatusCreated, "Product created", product)
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

// @Summary	Update a product
// @Tags	products
// @Accept	json
// @Produce	json
// @Param	id path int true "Product ID"
// @Param	body body request.UpdateProductRequest true "Product payload"
// @Success	200 {object} response.ProductAPIResponse
// @Failure	400 {object} response.ErrorResponse
// @Failure	500 {object} response.ErrorResponse
// @Router	/products/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var req request.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	pkgResponse.Success(c, http.StatusOK, "Product updated", product)
}

// @Summary	Delete a product
// @Tags	products
// @Produce	json
// @Param	id path int true "Product ID"
// @Success	200 {object} response.ErrorResponse
// @Failure	400 {object} response.ErrorResponse
// @Failure	500 {object} response.ErrorResponse
// @Router	/products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	pkgResponse.Success(c, http.StatusOK, "Product deleted", nil)
}
