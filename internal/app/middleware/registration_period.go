package middleware

import (
	"course-reg/internal/pkg/setting"
	"course-reg/internal/pkg/util"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CheckRegistrationPeriod checks if current time is within registration period
func CheckRegistrationPeriod(timeProvider util.TimeProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		now := timeProvider.Now()
		withinPeriod, err := setting.IsWithinRegistrationPeriod(now)
		if err != nil {
			log.Println("failed to check registration period:", err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "수강신청 기간 확인 실패"})
			return
		}
		if !withinPeriod {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "수강신청 기간이 아닙니다"})
			return
		}
		c.Next()
	}
}
