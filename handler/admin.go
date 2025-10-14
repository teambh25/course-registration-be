package handler

import (
	"course-reg/models"
	"course-reg/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(s *service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: s}
}

func (h *AdminHandler) RegisterStudents(c *gin.Context) {
	var students []models.Student

	if err := c.ShouldBindJSON(&students); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.adminService.RegisterStudents(students)

	c.JSON(http.StatusOK, gin.H{"message": "학생 등록 성공", "count": len(students)})
}

func (h *AdminHandler) CreateCourse(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "12134123412342134"})

}
