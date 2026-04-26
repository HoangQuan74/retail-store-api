package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoangquan/retail-store-api/pkg/model/request"
	pkgResponse "github.com/hoangquan/retail-store-api/pkg/response"
	"github.com/hoangquan/retail-store-api/services/api/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{service: authService}
}

// @Summary	Register a new user
// @Tags	auth
// @Accept	json
// @Produce	json
// @Param	body body request.RegisterRequest true "Register payload"
// @Success	201 {object} response.AuthAPIResponse
// @Failure	400 {object} response.ErrorResponse
// @Failure	409 {object} response.ErrorResponse
// @Router	/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		if err == service.ErrEmailAlreadyExists {
			pkgResponse.Error(c, http.StatusConflict, err.Error())
			return
		}
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	pkgResponse.Success(c, http.StatusCreated, "User registered", resp)
}

// @Summary	Login
// @Tags	auth
// @Accept	json
// @Produce	json
// @Param	body body request.LoginRequest true "Login payload"
// @Success	200 {object} response.AuthAPIResponse
// @Failure	400 {object} response.ErrorResponse
// @Failure	401 {object} response.ErrorResponse
// @Router	/auth/login [post]
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
