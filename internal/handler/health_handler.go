package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	pkgResponse "github.com/kainguyen/retail-store-api/pkg/response"
)

// @Summary	Health check
// @Tags	health
// @Produce	json
// @Success	200	{object}	response.CategoryAPIResponse
// @Router	/health [get]
func HealthCheck(c *gin.Context) {
	pkgResponse.Success(c, http.StatusOK, "OK", gin.H{
		"status": "healthy",
	})
}
