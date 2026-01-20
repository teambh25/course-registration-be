package session

import (
	"course-reg/internal/pkg/constant"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetSession(c *gin.Context) (constant.UserRole, uint, error) {
	session := sessions.Default(c)
	roleInt, ok := session.Get("role").(int)
	if !ok {
		return 0, 0, fmt.Errorf("get role failed")
	}

	userID, ok := session.Get("userID").(uint)
	if !ok {
		return 0, 0, fmt.Errorf("get user id failed")
	}

	return constant.UserRole(roleInt), userID, nil
}

func SetSession(c *gin.Context, role constant.UserRole, userID uint) error {
	session := sessions.Default(c)
	session.Set("role", int(role))
	session.Set("userID", userID)
	session.Options(sessions.Options{
		MaxAge:   60 * 60,
		Path:     "/",
		Domain:   "",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, // cf. 백엔드 서버와 프론트 서버의 도메인이 다르면 Domain과 SameStie 설정 필요
		// Secure:   isProd, // 운영 환경(HTTPS)일 때만 true
	})
	err := session.Save()
	return err
}

func DeleteSession(c *gin.Context) error {
	session := sessions.Default(c)
	session.Clear()
	err := session.Save()
	return err
}
