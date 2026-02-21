package session

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	cookieName      = "course_reg_session"
	maxAge          = 60 * 60 // 1 hour
	cleanupInterval = 30 * time.Minute
)

type sessionData struct {
	Role      UserRole
	UserID    uint
	ExpiresAt time.Time
}

var sessionStore sync.Map // map[string]sessionData

func newSessionID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func GetSession(c *gin.Context) (UserRole, uint, error) {
	id, err := c.Cookie(cookieName)
	if err != nil {
		return 0, 0, fmt.Errorf("no session cookie")
	}

	v, ok := sessionStore.Load(id)
	if !ok {
		return 0, 0, fmt.Errorf("session not found")
	}

	data := v.(sessionData)

	// lazy deletion
	if time.Now().After(data.ExpiresAt) {
		sessionStore.Delete(id)
		return 0, 0, fmt.Errorf("session expired")
	}

	return data.Role, data.UserID, nil
}

func CreateSession(c *gin.Context, role UserRole, userID uint) error {
	id, err := newSessionID()
	if err != nil {
		return err
	}

	sessionStore.Store(id, sessionData{
		Role:      role,
		UserID:    userID,
		ExpiresAt: time.Now().Add(maxAge * time.Second),
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     cookieName,
		Value:    id,
		MaxAge:   0,
		Path:     "/",
		HttpOnly: true,
		SameSite: cookieSameSite(),
		Secure:   cookieSecure(),
	})
	return nil
}

func DeleteSession(c *gin.Context) error {
	id, err := c.Cookie(cookieName)
	if err != nil {
		return nil // 이미 없으면 성공으로 처리?????????
	}

	sessionStore.Delete(id)

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		SameSite: cookieSameSite(),
		Secure:   cookieSecure(),
	})
	return nil
}

func StartCleanup() {
	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			sessionStore.Range(func(key, value any) bool {
				if now.After(value.(sessionData).ExpiresAt) {
					sessionStore.Delete(key)
				}
				return true
			})
		}
	}()
}

func cookieSameSite() http.SameSite {
	if gin.Mode() == gin.ReleaseMode {
		return http.SameSiteNoneMode
	}
	return http.SameSiteLaxMode
}

func cookieSecure() bool {
	return gin.Mode() == gin.ReleaseMode
}
