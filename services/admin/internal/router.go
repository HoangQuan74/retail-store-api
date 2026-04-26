package internal

import (
	"github.com/gin-gonic/gin"

	"github.com/hoangquan/retail-store-api/pkg/middleware"
	"github.com/hoangquan/retail-store-api/services/admin/internal/handler"
	"github.com/hoangquan/retail-store-api/services/admin/internal/service"
)

func NewRouter(deps *Dependencies) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())

	router.GET("/health", handler.HealthCheck)

	// Auth endpoints (public)
	authService := service.NewAuthService(deps.Queries, deps.JWTManager)
	authHandler := handler.NewAuthHandler(authService)
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/login", authHandler.Login)
	}

	// Admin-only endpoints (require JWT + admin role)
	admin := router.Group("/api/v1")
	admin.Use(middleware.Auth(deps.JWTManager))
	admin.Use(middleware.RequireRole("admin"))

	productService := service.NewProductService(deps.Queries, deps.Publisher)
	productHandler := handler.NewProductHandler(productService)
	products := admin.Group("/products")
	{
		products.POST("", productHandler.Create)
		products.GET("", productHandler.List)
		products.GET("/:id", productHandler.GetByID)
		products.PUT("/:id", productHandler.Update)
		products.DELETE("/:id", productHandler.Delete)
	}

	categoryService := service.NewCategoryService(deps.Queries)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	categories := admin.Group("/categories")
	{
		categories.POST("", categoryHandler.Create)
		categories.GET("", categoryHandler.List)
		categories.GET("/:id", categoryHandler.GetByID)
		categories.PUT("/:id", categoryHandler.Update)
		categories.DELETE("/:id", categoryHandler.Delete)
	}

	searchService := service.NewSearchService(deps.ESClient, deps.ProductIndex)
	searchHandler := handler.NewSearchHandler(searchService)
	admin.GET("/search/products", searchHandler.SearchProducts)

	return router
}
