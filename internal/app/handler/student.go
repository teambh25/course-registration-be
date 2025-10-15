package handler

import (
	"course-reg/internal/app/service"
)

type StudentHandler struct {
	studentService *service.StudentService
}

func NewStudentHandelr(s *service.StudentService) *StudentHandler {
	return &StudentHandler{studentService: s}
}
