package handler

import (
	"log"
	"net/http"
	"server/internal/models"
	"server/internal/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	service *service.InventoryService
}

func NewInventoryHandler(s *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: s}
}

func (h *InventoryHandler) CreateLog(c *gin.Context) {
	var input models.InventoryLog
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.service.CreateLog(c, input)
	if err != nil {
		log.Printf("Failed to create inventory log: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "log created"})
}

func (h *InventoryHandler) GetAllLogs(c *gin.Context) {
	var filter models.LogFilter

	if v := c.Query("updated"); v != "" {
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid value for 'updated', must be true or false"})
			return
		}
		filter.Updated = &parsed
	}

	if v := c.Query("product_id"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid value for 'product_id', must be an integer"})
			return
		}
		filter.ProductID = &parsed
	}

	if v := c.Query("date_from"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'date_from' format, use YYYY-MM-DD"})
			return
		}
		start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
		filter.DateFrom = &start
	}

	if v := c.Query("date_to"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'date_to' format, use YYYY-MM-DD"})
			return
		}
		end := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, time.UTC)
		filter.DateTo = &end
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

	result, err := h.service.GetAllLogs(c, filter, pagination)
	if err != nil {
		log.Printf("Failed to get inventory logs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *InventoryHandler) MarkLogsUpdated(c *gin.Context) {
	var body struct {
		IDs []int `json:"ids"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.MarkLogsUpdated(c, body.IDs); err != nil {
		log.Printf("Failed to mark logs as updated: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if len(body.IDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "all pending logs marked as updated"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "specified logs marked as updated", "ids": body.IDs})
	}
}