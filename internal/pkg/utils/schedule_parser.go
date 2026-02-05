package utils

import (
	"fmt"
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
// Format is fixed: "요일 HH:MM~HH:MM" (15 bytes: 한글 3 + 공백 1 + 시간 11)
func ParseSchedule(schedules string) ([]TimeSlot, error) {
	if schedules == "" {
		return nil, fmt.Errorf("empty schedule string")
	}

	var slots []TimeSlot
	parts := strings.Split(schedules, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		// "월 09:10~11:30" = 15 bytes
		if len(part) != 15 {
			return nil, fmt.Errorf("invalid schedule format: %q (length: %d, expected: 15)", part, len(part))
		}

		day := part[:3] // 한글 3바이트
		startHour, err := strconv.Atoi(part[4:6])
		if err != nil {
			return nil, fmt.Errorf("failed to parse start hour: %q", part[4:6])
		}
		startMin, err := strconv.Atoi(part[7:9])
		if err != nil {
			return nil, fmt.Errorf("failed to parse start minute: %q", part[7:9])
		}
		endHour, err := strconv.Atoi(part[10:12])
		if err != nil {
			return nil, fmt.Errorf("failed to parse end hour: %q", part[10:12])
		}
		endMin, err := strconv.Atoi(part[13:15])
		if err != nil {
			return nil, fmt.Errorf("failed to parse end minute: %q", part[13:15])
		}

		slots = append(slots, TimeSlot{
			Day:       day,
			StartHour: startHour,
			StartMin:  startMin,
			EndHour:   endHour,
			EndMin:    endMin,
		})
	}

	return slots, nil
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
func SchedulesConflict(schedule1, schedule2 string) (bool, error) {
	slots1, err := ParseSchedule(schedule1)
	if err != nil {
		return false, fmt.Errorf("schedule1: %w", err)
	}
	slots2, err := ParseSchedule(schedule2)
	if err != nil {
		return false, fmt.Errorf("schedule2: %w", err)
	}

	for _, s1 := range slots1 {
		for _, s2 := range slots2 {
			if HasConflict(s1, s2) {
				return true, nil
			}
		}
	}

	return false, nil
}
