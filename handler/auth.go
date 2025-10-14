package handler

import (
	"course-reg/service"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: s}
}

type user struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var u user

	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := h.authService.Check(u.Username, u.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Incorrect username"})
		return
	}

	if role != 0 {
		session := sessions.Default(c)
		session.Set("role", int(role))
		session.Options(sessions.Options{MaxAge: 60 * 60}) // 1 hour
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Log in!"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
	}
}

func (h *AuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusAccepted)
	c.JSON(http.StatusOK, gin.H{"message": "signed out"})
}
