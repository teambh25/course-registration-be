package constant

type CourseStatus string

const (
	CourseAvailable CourseStatus = "AVAILABLE" // 수강 신청 가능
	CourseWaitlist  CourseStatus = "WAITLIST"  // 대기자 신청 가능
	CourseFull      CourseStatus = "FULL"      // 정원 마감
)
