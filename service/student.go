package service

import (
	"course-reg/repository"
)

type StudentService struct {
	enrollRepo *repository.EnrollmentRepository
}

func NewStudentService(e *repository.EnrollmentRepository) *StudentService {
	return &StudentService{enrollRepo: e}
}
