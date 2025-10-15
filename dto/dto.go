package dto

type DeleteCourseRequset struct {
	CourseID uint `json:"course_id" binding:"required"`
}
