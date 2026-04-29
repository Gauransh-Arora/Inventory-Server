package handler

import (
	"net/http"
	"server/internal/models"
	"server/internal/service"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct{
	service *service.InventoryService
}

func NewInventoryHandler(s *service.InventoryService) *InventoryHandler{
	return &InventoryHandler{service: s}
}

func(h* InventoryHandler) CreateLog(c *gin.Context){
	var input models.InventoryLog
	if err := c.ShouldBindJSON(&input); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}
	err:=h.service.CreateLog(c, input)
	if err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message":"log created"})
}

func (h *InventoryHandler) GetAllLogs(c *gin.Context){
	logs,err:=h.service.GetAllLogs(c)
	if err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,logs)
}