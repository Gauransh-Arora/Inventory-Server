package middleware

import (
	"fmt"
	"net/http"
	"server/internal/repository"
	"server/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(repo *repository.AuthRepository) gin.HandlerFunc{
	return func(c *gin.Context){
		auth:=c.GetHeader("Authorization")
		if !strings.HasPrefix(auth,"Bearer "){
			c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{"error":"Bearer token required"})
			return
		}
		tokenStr := strings.TrimPrefix(auth,"Bearer ")
		token,err := jwt.Parse(tokenStr,func(t *jwt.Token)(interface{}, error){
			if _,ok := t.Method.(*jwt.SigningMethodRSA); !ok{
				return nil, fmt.Errorf("unexpected signing method")
			}
			return utils.GetVerifyKey(),nil
		})
		if err != nil || !token.Valid{
			c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{"error":"invalid or expired token"})
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		jti := claims["jti"].(string)
		if repo.IsDenyListed(c.Request.Context(), jti) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token has been revoked",
			})
			return
		}
		c.Set("userId", claims["userId"])
		c.Set("jti", jti)
		c.Set("exp", int64(claims["exp"].(float64)))
		c.Next()
	}
}
