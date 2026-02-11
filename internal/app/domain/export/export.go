package export

import (
	"course-reg/internal/app/repository"
	"course-reg/internal/pkg/file"
	"log"
	"sync"
)

const StaticCoursesFilePath = "static/courses.json"

var exportMu sync.Mutex

func ExportCoursesToJson(courseRepo repository.CourseRepositoryInterface) error {
	filePath := StaticCoursesFilePath
	exportMu.Lock()
	defer exportMu.Unlock()

	courses, err := courseRepo.FetchAllCourses()
	if err != nil {
		log.Println("[error] fetch all courses failed:", err.Error())
		return err
	}

	err = file.SaveJSON(filePath, courses)
	if err != nil {
		log.Println("[error] save JSON failed:", err.Error())
		return err
	}

	log.Printf("[info] Successfully exported %d courses to %s", len(courses), filePath)
	return nil
}
