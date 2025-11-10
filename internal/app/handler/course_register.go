package handler

import (
	"course-reg/internal/app/domain/dto"
	"course-reg/internal/app/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CourseRegHandler struct {
	courseRegService service.CourseRegServiceInterface
}

func NewCourseRegHandler(s service.CourseRegServiceInterface) *CourseRegHandler {
	return &CourseRegHandler{courseRegService: s}
}

func (h *CourseRegHandler) EnrollCourse(c *gin.Context) {
	studentID, ok := c.MustGet("studentID").(uint)
	if !ok {
		log.Println("enroll course failed: get student id failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "학생만 수강 신청이 가능합니다"})
		return
	}

	var req dto.EnrollCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("enroll course failed:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 수강 신청 요청"})
		return
	}

	success, message := h.courseRegService.Enroll(studentID, req.CourseID)

	c.JSON(http.StatusOK, gin.H{
		"success": success,
		"message": message,
		// "remaining_seats": allSeats,
	})
}

// func (h *CourseRegHandler) CancelEnrollment(c *gin.Context) {
// 	// todo: Get studentID from session
// 	studentID := uint(1)

// 	courseID, err := strconv.Atoi(c.Param("course_id"))
// 	if err != nil {
// 		log.Println("cancel enrollment failed:", err.Error())
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 강의 ID"})
// 		return
// 	}

// 	success, message, allSeats := h.courseRegService.CancelEnrollment(studentID, uint(courseID))
// 	c.JSON(http.StatusOK, gin.H{
// 		"success":         success,
// 		"message":         message,
// 		"remaining_seats": allSeats,
// 	})
// }

// ForceEnrollCourse
// ForceCancelCourse
// GetStatus
