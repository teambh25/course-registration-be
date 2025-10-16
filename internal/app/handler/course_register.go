package handler

import (
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

func (h *CourseRegHandler) GetAllCourses(c *gin.Context) {
	courses, err := h.courseRegService.GetAllCourses()
	if err != nil {
		log.Println("get courses failed:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "전체 강의 가져오기 실패"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"course_id": courses})
}
