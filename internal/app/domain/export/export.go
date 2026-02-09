package export

import (
	"course-reg/internal/app/repository"
	"course-reg/internal/pkg/file"
	"course-reg/internal/pkg/setting"
	"log"
	"sync"
)

var exportMu sync.Mutex

func ExportCoursesToJson(courseRepo repository.CourseRepositoryInterface) error {
	exportMu.Lock()
	defer exportMu.Unlock()

	courses, err := courseRepo.FetchAllCourses()
	if err != nil {
		log.Println("[error] fetch all courses failed:", err.Error())
		return err
	}

	err = file.SaveJSON(setting.AppSetting.StaticCoursesFilePath, courses)
	if err != nil {
		log.Println("[error] save JSON failed:", err.Error())
		return err
	}

	log.Printf("[info] Successfully exported %d courses to %s", len(courses), setting.AppSetting.StaticCoursesFilePath)
	return nil
}
