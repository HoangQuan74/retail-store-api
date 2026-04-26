package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoangquan/retail-store-api/pkg/model/request"
	pkgResponse "github.com/hoangquan/retail-store-api/pkg/response"
	"github.com/hoangquan/retail-store-api/services/admin/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{service: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			pkgResponse.Error(c, http.StatusUnauthorized, err.Error())
			return
		}
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	pkgResponse.Success(c, http.StatusOK, "Login successful", resp)
}
