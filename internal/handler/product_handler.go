package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/kainguyen/retail-store-api/db/sqlc"
	"github.com/kainguyen/retail-store-api/internal/model/request"
	"github.com/kainguyen/retail-store-api/internal/service"
	appnats "github.com/kainguyen/retail-store-api/pkg/nats"
	pkgResponse "github.com/kainguyen/retail-store-api/pkg/response"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(queries *db.Queries, publisher *appnats.Publisher) *ProductHandler {
	return &ProductHandler{service: service.NewProductService(queries, publisher)}
}

func (h *ProductHandler) RegisterRoutes(router *gin.Engine) {
	products := router.Group("/api/v1/products")
	{
		products.POST("", h.Create)
		products.GET("", h.List)
		products.GET("/:id", h.GetByID)
		products.PUT("/:id", h.Update)
		products.DELETE("/:id", h.Delete)
	}
}

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
