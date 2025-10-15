package handler

import (
	"course-reg/service"
)

type StudentHandler struct {
	studentService *service.StudentService
}

func NewStudentHandelr(s *service.StudentService) *StudentHandler {
	return &StudentHandler{studentService: s}
}
