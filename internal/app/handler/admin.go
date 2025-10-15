package handler

import (
	"course-reg/internal/app/dto"
	"course-reg/internal/app/models"
	"course-reg/internal/app/service"
	"log"
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
		// todo: json vadidation
		log.Println("register students failed:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 학생 형식"})
		return
	}

	if err := h.adminService.RegisterStudents(students); err != nil {
		// todo: 중복된 학생 처리
		c.JSON(http.StatusBadRequest, gin.H{"error": "학생 등록 실패"})
		return
	}

	c.Status(http.StatusOK)
}

func (h *AdminHandler) ResetStudents(c *gin.Context) {
	if err := h.adminService.ResetStudents(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "학생 삭제 실패, 개발자 호출 필요!"})
	}

	c.Status(http.StatusOK)
}

func (h *AdminHandler) CreateCourse(c *gin.Context) {
	var course = &models.Course{}

	if err := c.ShouldBindJSON(course); err != nil {
		log.Println("create course failed:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 강의 형식"})
		return
	}

	courseID, err := h.adminService.CreateCourse(course)
	if err != nil {
		// todo: 중복된 강의 처리
		c.JSON(http.StatusBadRequest, gin.H{"error": "강의 등록 실패"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"course_id": courseID})
}

func (h *AdminHandler) DeleteCourse(c *gin.Context) {
	var req dto.DeleteCourseRequset

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("delete course failed:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 코스 id"})
		return
	}

	if err := h.adminService.DeleteCourse(req.CourseID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "강의 삭제 실패"})
	}

	c.Status(http.StatusOK)
}
