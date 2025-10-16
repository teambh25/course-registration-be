package session

import (
	"course-reg/internal/pkg/constant"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetSession(c *gin.Context) (constant.UserRole, error) {
	session := sessions.Default(c)
	roleInt, ok := session.Get("role").(int)
	if !ok {
		return 0, fmt.Errorf("get session failed")
	}
	return constant.UserRole(roleInt), nil
}

func DeleteSession(c *gin.Context) error {
	session := sessions.Default(c)
	session.Clear()
	err := session.Save()
	return err
}

func SetSession(c *gin.Context, role constant.UserRole) error {
	session := sessions.Default(c)
	session.Set("role", int(role))                     // json encoding 해서 저장하는게 일반적?
	session.Options(sessions.Options{MaxAge: 60 * 60}) // 1 hour
	err := session.Save()
	return err
}
