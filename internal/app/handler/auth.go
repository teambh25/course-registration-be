package handler

import (
	"course-reg/internal/app/service"
	"course-reg/internal/pkg/session"
	"log"
	"net/http"

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
	if err != nil || role == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "잘못된 ID/PW"})
		return
	}

	err = session.SetSession(c, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"role": role.String()})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	err := session.DeleteSession(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *AuthHandler) Check(c *gin.Context) {
	role, err := session.GetSession(c)
	if err != nil {
		log.Println("auth check failed:", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "세션 만료"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"role": role.String()})
}
