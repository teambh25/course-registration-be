package handler

import (
	"course-reg/internal/app/domain/dto"
	"course-reg/internal/app/domain/e"
	"course-reg/internal/app/service"
	"errors"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "학생만 수강 신청이 가능합니다"})
		return
	}

	var req dto.EnrollCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("enroll course :", err.Error()) // If this occurs, check the client-side request
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 수강 신청 요청"})
		return
	}

	err := h.courseRegService.Enroll(studentID, req.CourseID)
	if err != nil {
		status, msg := enrollErrToResponse(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "수강신청 성공"})
}

func enrollErrToResponse(err error) (int, string) {
	switch {
	case errors.Is(err, e.ErrCourseNotFound):
		return http.StatusNotFound, "존재하지 않는 강의입니다"
	case errors.Is(err, e.ErrStudentNotFound):
		return http.StatusNotFound, "존재하지 않는 학생입니다"
	case errors.Is(err, e.ErrTimeConflict):
		return http.StatusConflict, "시간이 겹치는 강의가 있습니다"
	case errors.Is(err, e.ErrAlreadyEnrolled):
		return http.StatusConflict, "이미 신청한 강의입니다"
	case errors.Is(err, e.ErrCourseFull):
		return http.StatusConflict, "정원이 초과되었습니다"
	case errors.Is(err, e.ErrInvalidRegistrationPeriod):
		return http.StatusForbidden, "수강신청 기간이 아닙니다"
	default:
		log.Println("enroll unexpected error:", err.Error())
		return http.StatusInternalServerError, "알 수 없는 오류가 발생했습니다"
	}
}

func (h *CourseRegHandler) GetAllCourseStatus(c *gin.Context) {
	result, err := h.courseRegService.GetAllCourseStatus()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// func (h *CourseRegHandler) CancelEnrollment(c *gin.Context) {
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
