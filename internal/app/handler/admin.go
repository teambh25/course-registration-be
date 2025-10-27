package handler

import (
	"course-reg/internal/app/domain/dto"
	"course-reg/internal/app/models"
	"course-reg/internal/app/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminService service.AdminServiceInterface
}

func NewAdminHandler(s service.AdminServiceInterface) *AdminHandler {
	return &AdminHandler{adminService: s}
}

func (h *AdminHandler) RegisterStudents(c *gin.Context) {
	var students []models.Student

	if err := c.ShouldBindJSON(&students); err != nil {
		// todo: json vadidation
		log.Println("register students failed:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 학생 리스트"})
		return
	}

	if err := h.adminService.RegisterStudents(students); err != nil {
		// todo: 중복된 학생 처리
		c.JSON(http.StatusBadRequest, gin.H{"error": "학생 리스트 등록 실패"})
		return
	}

	c.Status(http.StatusOK)
}

func (h *AdminHandler) ResetStudents(c *gin.Context) {
	if err := h.adminService.ResetStudents(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "학생 리스트 삭제 실패, 개발자 호출 필요!"})
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
	course_id, err := strconv.Atoi(c.Param("course_id"))
	if err != nil {
		log.Println("delete course failed:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 강의 id"})
		return
	}

	if err := h.adminService.DeleteCourse(uint(course_id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "강의 삭제 실패"})
	}

	c.Status(http.StatusOK)
}

func (h *AdminHandler) RegisterCourses(c *gin.Context) {
	var courses []models.Course

	if err := c.ShouldBindJSON(&courses); err != nil {
		log.Println("register courses failed:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 강의 리스트"})
		return
	}

	if err := h.adminService.RegisterCourses(courses); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "강의 리스트 등록 실패"})
		return
	}

	c.Status(http.StatusOK)
}

func (h *AdminHandler) ResetCourses(c *gin.Context) {
	if err := h.adminService.ResetCourses(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "강의 리스트 삭제 실패, 개발자 호출 필요!"})
	}

	c.Status(http.StatusOK)
}

func (h *AdminHandler) SetRegistrationPeriod(c *gin.Context) {
	var req dto.SetRegistrationPeriodRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("set registration period failed:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청 형식"})
		return
	}

	if err := h.adminService.SetRegistrationPeriod(req.StartTime, req.EndTime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ""})
		return
	}

	c.Status(http.StatusOK)
}

func (h *AdminHandler) GetRegistrationPeriod(c *gin.Context) {
	startTime, endTime, err := h.adminService.GetRegistrationPeriod()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "수강 신청 기간 조회 실패"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"start_time": startTime,
		"end_time":   endTime,
	})
}
