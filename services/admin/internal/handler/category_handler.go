package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hoangquan/retail-store-api/pkg/model/request"
	pkgResponse "github.com/hoangquan/retail-store-api/pkg/response"
	"github.com/hoangquan/retail-store-api/services/admin/internal/service"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: svc}
}

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

func (h *CategoryHandler) List(c *gin.Context) {
	categories, err := h.service.List(c.Request.Context())
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	pkgResponse.Success(c, http.StatusOK, "Success", categories)
}

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
