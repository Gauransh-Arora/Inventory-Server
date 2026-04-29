package routes

import (
	"server/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, inv *handler.InventoryHandler, prod *handler.ProductHandler){
	//inventory logs
	r.POST("/logs",inv.CreateLog)
	r.GET("/logs",inv.GetAllLogs)

	//products
	r.POST("/products",prod.CreateProduct)
	r.GET("/products",prod.GetAllProducts)
	r.GET("/products/:barcode",prod.GetProductByBarcode)
}