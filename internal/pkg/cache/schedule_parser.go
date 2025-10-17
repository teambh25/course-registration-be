package cache

import (
	"regexp"
	"strconv"
	"strings"
)

// TimeSlot represents a single time slot (day + start-end time)
type TimeSlot struct {
	Day       string // "월", "화", "수", "목", "금", "토", "일"
	StartHour int    //
	StartMin  int    //
	EndHour   int    //
	EndMin    int    //
}

// ParseSchedule parses schedule string like "월 09:10~11:30, 수 17:10~19:20"
func ParseSchedule(schedules string) []TimeSlot {
	var slots []TimeSlot

	// Split by comma
	parts := strings.Split(schedules, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Pattern: "월 09:10~11:30"
		re := regexp.MustCompile(`([월화수목금토일])\s+(\d{2}):(\d{2})~(\d{2}):(\d{2})`)
		matches := re.FindStringSubmatch(part)

		if len(matches) != 6 {
			continue
		}

		day := matches[1]
		startHour, _ := strconv.Atoi(matches[2])
		startMin, _ := strconv.Atoi(matches[3])
		endHour, _ := strconv.Atoi(matches[4])
		endMin, _ := strconv.Atoi(matches[5])

		slots = append(slots, TimeSlot{
			Day:       day,
			StartHour: startHour,
			StartMin:  startMin,
			EndHour:   endHour,
			EndMin:    endMin,
		})
	}

	return slots
}

// HasConflict checks if two time slots conflict
func HasConflict(slot1, slot2 TimeSlot) bool {
	// Different days - no conflict
	if slot1.Day != slot2.Day {
		return false
	}

	// Same day - check time overlap
	// Convert to minutes for easier comparison
	start1 := slot1.StartHour*60 + slot1.StartMin
	end1 := slot1.EndHour*60 + slot1.EndMin
	start2 := slot2.StartHour*60 + slot2.StartMin
	end2 := slot2.EndHour*60 + slot2.EndMin

	// Conflict if: start1 < end2 AND start2 < end1
	return start1 < end2 && start2 < end1
}

// SchedulesConflict checks if two schedule strings conflict
func SchedulesConflict(schedule1, schedule2 string) bool {
	slots1 := ParseSchedule(schedule1)
	slots2 := ParseSchedule(schedule2)

	for _, s1 := range slots1 {
		for _, s2 := range slots2 {
			if HasConflict(s1, s2) {
				return true
			}
		}
	}

	return false
}
