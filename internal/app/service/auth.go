package service

import (
	"course-reg/internal/app/repository"
	"course-reg/internal/pkg/constant"
	"course-reg/internal/pkg/setting"
)

const (
	ADMIN   int = 0
	STUDENT int = 1
)

type AuthService struct {
	studentRepo repository.StudentRepositoryInterface
}

func NewAuthService(s repository.StudentRepositoryInterface) *AuthService {
	return &AuthService{studentRepo: s}
}

func (a *AuthService) Check(username string, password string) (constant.UserRole, uint, error) {
	var role constant.UserRole
	var pw string
	var userID uint
	var err error

	if is_admin := setting.SecretSetting.AdminID == username && setting.SecretSetting.AdminPW == password; is_admin {
		role = constant.RoleAdmin
	} else if userID, pw, err = a.studentRepo.FetchPassword(username); err == nil && pw == password {
		role = constant.RoleStudent
	}
	return role, userID, err
}
