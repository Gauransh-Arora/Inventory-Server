package handler

import (
	"log"
	"server/internal/models"
	"server/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	Service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{Service: s}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid input",
		})
		return
	}
	id, err := h.Service.Register(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(409, gin.H{
			"error": "user already exists",
		})
		return
	}
	c.JSON(201, gin.H{
		"id":      id,
		"username":   req.Username,
		"message": "User created successfully",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Credentials required"})
		return
	}
	resp, err := h.Service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, resp)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Refresh token required"})
		return
	}

	resp, err := h.Service.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if err.Error() == "invalid refresh token" {
			c.JSON(401, gin.H{"error": "invalid refresh token"})
		} else {
			log.Printf("Token refresh failed: %v", err)
			c.JSON(500, gin.H{"error": "internal server error"})
		}
		return
	}
	c.JSON(200, resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userIDStr, _ := c.Get("userId")
	userID, _ := uuid.Parse(userIDStr.(string))
	jti, _ := c.Get("jti")
	exp, _ := c.Get("exp")

	err := h.Service.Logout(c.Request.Context(), jti.(string), exp.(int64), userID)
	if err != nil {
		log.Printf("Logout failed for user %s: %v", userIDStr, err)
		c.JSON(500, gin.H{"error": "failed to logout"})
		return
	}
	c.JSON(200, gin.H{"message": "Logged out successfully"})
}