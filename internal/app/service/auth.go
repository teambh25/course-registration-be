package service

import (
	"course-reg/internal/app/repository"
	"course-reg/internal/pkg/crypto"
	"course-reg/internal/pkg/session"
	"log"
)

type AuthService struct {
	studentRepo repository.StudentRepositoryInterface
	adminID     string
	adminPW     string
	pepper      string
}

func NewAuthService(s repository.StudentRepositoryInterface, adminID, adminPW, pepper string) *AuthService {
	return &AuthService{studentRepo: s, adminID: adminID, adminPW: adminPW, pepper: pepper}
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
		} else if crypto.VerifyPassword(pw, password, a.pepper) {
			role = session.RoleStudent
		}
	}
	return role, userID, err
}
