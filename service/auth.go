package service

import (
	"course-reg/pkg/constant"
	"course-reg/pkg/setting"
	"course-reg/repository"
	"fmt"
)

const (
	ADMIN   int = 0
	STUDENT int = 1
)

type AuthService struct {
	studentRepo *repository.StudentRepository
}

func NewAuthService(s *repository.StudentRepository) *AuthService {
	return &AuthService{studentRepo: s}
}

func (a *AuthService) Check(username string, password string) (constant.UserRole, error) {
	var role constant.UserRole
	var err error
	var pw string

	if is_admin := setting.SecretSetting.AdminID == username && setting.SecretSetting.AdminPW == password; is_admin {
		role = constant.RoleAdmin
	} else if pw, err = a.studentRepo.GetPassword(username); err == nil && pw == password {
		role = constant.RoleStudent
	}

	fmt.Println("[check]", username, password, pw)
	fmt.Println("[check]", role, err)
	return role, err
}
