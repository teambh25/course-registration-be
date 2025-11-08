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
