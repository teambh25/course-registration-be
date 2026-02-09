package session

type UserRole int // todo: string으로 바꾸기?

const (
	RoleAdmin UserRole = iota + 1
	RoleStudent
	// RoleInstructor
)

func (r UserRole) String() string {
	switch r {
	case RoleAdmin:
		return "admin"
	case RoleStudent:
		return "student"
	// case RoleInstructor:
	// 	return "instructor"
	default:
		return "unknown"
	}
}
