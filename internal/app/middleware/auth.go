package middleware

import (
	"course-reg/internal/pkg/constant"
	"course-reg/internal/pkg/session"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, _, err := session.GetSession(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		}
		c.Next()
	}
}

func AuthAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _, err := session.GetSession(c)
		if err != nil || role != constant.RoleAdmin {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		}
		c.Next()
	}
}

func AuthStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, userID, err := session.GetSession(c)
		if err != nil || role != constant.RoleStudent {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		}
		c.Set("studentID", userID)
		c.Next()
	}
}
