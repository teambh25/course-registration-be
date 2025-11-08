package session

import (
	"course-reg/internal/pkg/constant"
	"fmt"

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
	session.Options(sessions.Options{MaxAge: 60 * 60}) // 1 hour
	err := session.Save()
	return err
}

func DeleteSession(c *gin.Context) error {
	session := sessions.Default(c)
	session.Clear()
	err := session.Save()
	return err
}
