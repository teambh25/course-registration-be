package static

import (
	"course-reg/internal/app/repository"
	"course-reg/internal/pkg/constant"
	"course-reg/internal/pkg/file"
	"log"
)

func ExportCoursesToJson(courseRepo repository.CourseRepositoryInterface) error {
	courses, err := courseRepo.FetchAllCourses()
	if err != nil {
		log.Println("fetch all courses failed:", err.Error())
		return err
	}

	err = file.SaveJSON(constant.StaticCoursesFilePath, courses)
	if err != nil {
		log.Println("save JSON failed:", err.Error())
		return err
	}

	log.Printf("Successfully exported %d courses to %s", len(courses), constant.StaticCoursesFilePath)
	return nil
}
