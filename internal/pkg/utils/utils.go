package utils

import "time"

// TimeProvider provides current time (for dependency injection and testing)
type TimeProvider interface {
	Now() time.Time
}

// KoreaTimeProvider returns current time in Korea timezone
type KoreaTimeProvider struct{}

func NewKoreaTimeProvider() *KoreaTimeProvider {
	return &KoreaTimeProvider{}
}

func (p *KoreaTimeProvider) Now() time.Time {
	loc, _ := time.LoadLocation("Asia/Seoul")
	return time.Now().In(loc)
}

func StringToTime(timeStr string) (time.Time, error) {
	return time.Parse("2006-01-02-15-04", timeStr)
}
