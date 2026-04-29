package handler

import (
	"net/http"
	"server/internal/models"
	"server/internal/service"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct{
	service *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler{
	return &ProductHandler{service: s}
}

//Post /Products
func(h *ProductHandler) CreateProduct(c *gin.Context){
	var input models.Product

	if err:=c.ShouldBindJSON(&input); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error":err.Error(),
		})
		return
	}

	err:=h.service.CreateProduct(c,input)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated,gin.H{
		"message":"product created",
	})
}

func(h *ProductHandler) GetAllProducts(c *gin.Context){
	products,err := h.service.GetAllProducts(c)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, products)
}

func(h *ProductHandler) GetProductByBarcode(c *gin.Context){
	barcode:=c.Param("barcode")
	products,err:=h.service.GetProductByBarcode(c,barcode)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":err.Error(),
		})
		return
	}
	if len(products) == 0{
		c.JSON(http.StatusNotFound,gin.H{
			"message":"no products found",
		})
		return
	}
	c.JSON(http.StatusOK, products)
}