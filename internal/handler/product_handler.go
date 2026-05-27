package handler

import (
	"log"
	"net/http"
	"server/internal/models"
	"server/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var input models.Product

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateProduct(c, input); err != nil {
		log.Printf("Failed to create product: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "product created"})
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	var icode *int

	if v := c.Query("icode"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid value for 'icode', must be an integer"})
			return
		}
		icode = &parsed
	}

	page := 1
	pageSize := 20

	if v := c.Query("page"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil || parsed < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid value for 'page', must be a positive integer"})
			return
		}
		page = parsed
	}

	if v := c.Query("page_size"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil || parsed < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid value for 'page_size', must be a positive integer"})
			return
		}
		if parsed > 100 {
			parsed = 100
		}
		pageSize = parsed
	}

	pagination := models.Pagination{Page: page, PageSize: pageSize}

	result, err := h.service.GetAllProducts(c, icode, pagination)
	if err != nil {
		log.Printf("Failed to get products: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *ProductHandler) GetProductByBarcode(c *gin.Context) {
	barcode := c.Param("barcode")
	products, err := h.service.GetProductByBarcode(c, barcode)
	if err != nil {
		log.Printf("Failed to get product by barcode %s: %v", barcode, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no products found"})
		return
	}
	c.JSON(http.StatusOK, products)
}
