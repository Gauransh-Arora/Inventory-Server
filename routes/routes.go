package routes

import (
	"server/internal/handler"
	"server/internal/middleware"
	"server/internal/repository"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, inv *handler.InventoryHandler, prod *handler.ProductHandler, auth *handler.AuthHandler, authRepo *repository.AuthRepository) {
	v1 := r.Group("/api/v1/auth")
	{
		v1.POST("/register", auth.Register)
		v1.POST("/login", auth.Login)
		v1.POST("/refresh", auth.Refresh)
	}

	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(authRepo))
	{
		protected.POST("/logout", auth.Logout)

		// inventory logs
		protected.POST("/logs", inv.CreateLog)
		protected.GET("/logs", inv.GetAllLogs)
		protected.PATCH("/logs", inv.MarkLogsUpdated)

		// products
		protected.POST("/products", prod.CreateProduct)
		protected.GET("/products", prod.GetAllProducts)
		protected.GET("/products/:barcode", prod.GetProductByBarcode)
		protected.PATCH("/products/:id", prod.UpdateProduct)
		protected.DELETE("/products", prod.DeleteProducts)
	}
}