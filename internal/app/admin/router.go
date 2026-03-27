package admin

import (
	"github.com/gin-gonic/gin"

	"github.com/hoangquan/retail-store-api/internal/app"
	"github.com/hoangquan/retail-store-api/internal/handler"
	"github.com/hoangquan/retail-store-api/internal/service"
	"github.com/hoangquan/retail-store-api/pkg/middleware"
)

func NewRouter(ctx *app.AppContext) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())

	router.GET("/health", handler.HealthCheck)

	// Auth endpoints (public)
	authService := service.NewAuthService(ctx.Queries, ctx.JWTManager)
	authHandler := handler.NewAuthHandler(authService)
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/login", authHandler.Login)
	}

	// Admin-only endpoints (require JWT + admin role)
	admin := router.Group("/api/v1")
	admin.Use(middleware.Auth(ctx.JWTManager))
	admin.Use(middleware.RequireRole("admin"))

	productHandler := handler.NewProductHandler(ctx)
	products := admin.Group("/products")
	{
		products.POST("", productHandler.Create)
		products.GET("", productHandler.List)
		products.GET("/:id", productHandler.GetByID)
		products.PUT("/:id", productHandler.Update)
		products.DELETE("/:id", productHandler.Delete)
	}

	categoryHandler := handler.NewCategoryHandler(ctx)
	categories := admin.Group("/categories")
	{
		categories.POST("", categoryHandler.Create)
		categories.GET("", categoryHandler.List)
		categories.GET("/:id", categoryHandler.GetByID)
		categories.PUT("/:id", categoryHandler.Update)
		categories.DELETE("/:id", categoryHandler.Delete)
	}

	searchHandler := handler.NewSearchHandler(ctx)
	admin.GET("/search/products", searchHandler.SearchProducts)

	return router
}
