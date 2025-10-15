package middleware

import (
	"course-reg/pkg/constant"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		role := session.Get("role")
		if x, ok := role.(int); !ok || constant.UserRole(x) != constant.RoleAdmin {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		}
		c.Next() // Request 이전
	}
}

func AuthStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		role := session.Get("role")
		if role == nil { // || role != constant.RoleStudent
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		}
		c.Next() // Request 이전
	}
}
