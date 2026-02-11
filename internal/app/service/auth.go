package service

import (
	"course-reg/internal/app/repository"
	"course-reg/internal/pkg/session"
	"log"
)

type AuthService struct {
	studentRepo repository.StudentRepositoryInterface
	adminID     string
	adminPW     string
}

func NewAuthService(s repository.StudentRepositoryInterface, adminID, adminPW string) *AuthService {
	return &AuthService{studentRepo: s, adminID: adminID, adminPW: adminPW}
}

func (a *AuthService) Check(username string, password string) (session.UserRole, uint, error) {
	var role session.UserRole
	var pw string
	var userID uint
	var err error

	if a.adminID == username && a.adminPW == password {
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
