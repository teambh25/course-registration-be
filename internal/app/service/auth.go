package service

import (
	"course-reg/internal/app/repository"
	"course-reg/internal/pkg/session"
	"course-reg/internal/pkg/setting"
	"log"
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

func (a *AuthService) Check(username string, password string) (session.UserRole, uint, error) {
	var role session.UserRole
	var pw string
	var userID uint
	var err error

	if is_admin := setting.SecretSetting.AdminID == username && setting.SecretSetting.AdminPW == password; is_admin {
		role = session.RoleAdmin
	} else {
		userID, pw, err = a.studentRepo.FetchPassword(username)
		if err != nil {
			log.Println("[error] fetch password failed", err.Error())
		} else if pw == password {
			role = session.RoleStudent
		}
	}
	return role, userID, err
}
