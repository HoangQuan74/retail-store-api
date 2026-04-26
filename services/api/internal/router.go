package internal

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/hoangquan/retail-store-api/docs"
	"github.com/hoangquan/retail-store-api/pkg/middleware"
	"github.com/hoangquan/retail-store-api/services/api/internal/handler"
	"github.com/hoangquan/retail-store-api/services/api/internal/service"
)

func NewRouter(deps *Dependencies) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())

	router.GET("/health", handler.HealthCheck)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth endpoints (public)
	authService := service.NewAuthService(deps.Queries, deps.JWTManager)
	authHandler := handler.NewAuthHandler(authService)
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	// Public read-only endpoints
	productService := service.NewProductService(deps.Queries)
	productHandler := handler.NewProductHandler(productService)
	products := router.Group("/api/v1/products")
	{
		products.GET("", productHandler.List)
		products.GET("/:id", productHandler.GetByID)
	}

	categoryService := service.NewCategoryService(deps.Queries)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	categories := router.Group("/api/v1/categories")
	{
		categories.GET("", categoryHandler.List)
		categories.GET("/:id", categoryHandler.GetByID)
	}

	searchService := service.NewSearchService(deps.ESClient, deps.ProductIndex)
	searchHandler := handler.NewSearchHandler(searchService)
	router.GET("/api/v1/search/products", searchHandler.SearchProducts)

	return router
}
