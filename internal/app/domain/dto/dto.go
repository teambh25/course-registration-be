package dto

type SetRegistrationPeriodRequest struct {
	StartTime string `json:"start_time" binding:"required"` // "2025-01-20-09-00"
	EndTime   string `json:"end_time" binding:"required"`   // "2025-01-25-18-00"
}

type EnrollCourseRequest struct {
	CourseID uint `json:"course_id" binding:"required"`
}
