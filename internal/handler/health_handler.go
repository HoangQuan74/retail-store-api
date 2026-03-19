package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	pkgResponse "github.com/kainguyen/retail-store-api/pkg/response"
)

func HealthCheck(c *gin.Context) {
	pkgResponse.Success(c, http.StatusOK, "OK", gin.H{
		"status": "healthy",
	})
}
