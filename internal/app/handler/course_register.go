package handler

import (
	"course-reg/internal/app/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CourseRegHandler struct {
	courseRegService service.CourseRegServiceInterface
}

func NewCourseRegHandler(s service.CourseRegServiceInterface) *CourseRegHandler {
	return &CourseRegHandler{courseRegService: s}
}

func (h *CourseRegHandler) GetAllCourses(c *gin.Context) {
	courses, err := h.courseRegService.GetAllCourses()
	if err != nil {
		log.Println("get courses failed:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "전체 강의 가져오기 실패"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"courses": courses})
}

func (h *CourseRegHandler) EnrollCourse(c *gin.Context) {
	// TODO: Get studentID from session
	studentID := uint(1) // placeholder

	courseID, err := strconv.Atoi(c.Param("course_id"))
	if err != nil {
		log.Println("enroll course failed:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 강의 ID"})
		return
	}

	success, message, allSeats, waitlistPos := h.courseRegService.Enroll(studentID, uint(courseID))

	c.JSON(http.StatusOK, gin.H{
		"success":           success,
		"message":           message,
		"remaining_seats":   allSeats,
		"waitlist_position": waitlistPos,
	})
}

func (h *CourseRegHandler) CancelEnrollment(c *gin.Context) {
	// TODO: Get studentID from session
	studentID := uint(1) // placeholder

	courseID, err := strconv.Atoi(c.Param("course_id"))
	if err != nil {
		log.Println("cancel enrollment failed:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 강의 ID"})
		return
	}

	success, message, allSeats := h.courseRegService.CancelEnrollment(studentID, uint(courseID))

	c.JSON(http.StatusOK, gin.H{
		"success":         success,
		"message":         message,
		"remaining_seats": allSeats,
	})
}
