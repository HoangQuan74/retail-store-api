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

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(ctx *app.AppContext, router *gin.Engine) *CategoryHandler {
	h := &CategoryHandler{service: service.NewCategoryService(ctx.Queries)}
	categories := router.Group("/api/v1/categories")
	{
		categories.POST("", h.Create)
		categories.GET("", h.List)
		categories.GET("/:id", h.GetByID)
		categories.PUT("/:id", h.Update)
		categories.DELETE("/:id", h.Delete)
	}
	return h
}

// @Summary	Create a category
// @Tags	categories
// @Accept	json
// @Produce	json
// @Param	body body request.CreateCategoryRequest true "Category payload"
// @Success	201 {object} response.CategoryAPIResponse
// @Failure	400 {object} response.ErrorResponse
// @Failure	500 {object} response.ErrorResponse
// @Router	/categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	var req request.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	pkgResponse.Success(c, http.StatusCreated, "Category created", category)
}

// swag-list
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

// @Summary	Update a category
// @Tags	categories
// @Accept	json
// @Produce	json
// @Param	id path int true "Category ID"
// @Param	body body request.UpdateCategoryRequest true "Category payload"
// @Success	200 {object} response.CategoryAPIResponse
// @Failure	400 {object} response.ErrorResponse
// @Failure	500 {object} response.ErrorResponse
// @Router	/categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var req request.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	pkgResponse.Success(c, http.StatusOK, "Category updated", category)
}

// @Summary	Delete a category
// @Tags	categories
// @Produce	json
// @Param	id path int true "Category ID"
// @Success	200 {object} response.ErrorResponse
// @Failure	400 {object} response.ErrorResponse
// @Failure	500 {object} response.ErrorResponse
// @Router	/categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	pkgResponse.Success(c, http.StatusOK, "Category deleted", nil)
}
