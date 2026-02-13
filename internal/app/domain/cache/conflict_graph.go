package cache

import (
	"course-reg/internal/app/models"
	"fmt"
	"strconv"
	"strings"
)

type ConflictGraph map[uint]map[uint]bool

// courseTime represents a single time slot (day + start-end time)
type courseTime struct {
	Day       string
	StartHour int
	StartMin  int
	EndHour   int
	EndMin    int
}

func BuildConflictGraph(courses []models.Course) (ConflictGraph, error) {
	graph := make(ConflictGraph)

	// Parse all schedules upfront for validation and caching
	parsed := make(map[uint][]courseTime)
	for _, c := range courses {
		slots, err := parseCourseTime(c.Schedules)
		if err != nil {
			return nil, fmt.Errorf("failed to parse schedule for course %d: %w", c.ID, err)
		}
		parsed[c.ID] = slots
		graph[c.ID] = make(map[uint]bool)
	}

	// Build conflict graph (only iterate i < j pairs)
	for i, c1 := range courses {
		for j := i + 1; j < len(courses); j++ {
			c2 := courses[j]
			if isScheduleConflict(parsed[c1.ID], parsed[c2.ID]) {
				graph[c1.ID][c2.ID] = true
				graph[c2.ID][c1.ID] = true
			}
		}
	}
	return graph, nil
}

// parseCourseTime parses schedule string like "월 09:10~11:30, 수 17:10~19:20"
// Format is fixed: "요일 HH:MM~HH:MM" (15 bytes: 3 Korean + 1 space + 11 time)
func parseCourseTime(schedules string) ([]courseTime, error) {
	if schedules == "" {
		return nil, fmt.Errorf("empty schedule string")
	}

	var slots []courseTime
	parts := strings.Split(schedules, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		// "월 09:10~11:30" = 15 bytes
		if len(part) != 15 {
			return nil, fmt.Errorf("invalid schedule format: %q (length: %d, expected: 15)", part, len(part))
		}

		day := part[:3] // 3 bytes for Korean char
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

		slots = append(slots, courseTime{
			Day:       day,
			StartHour: startHour,
			StartMin:  startMin,
			EndHour:   endHour,
			EndMin:    endMin,
		})
	}

	return slots, nil
}

// isScheduleConflict checks if any time slots between two slices conflict
func isScheduleConflict(schedule1, schedule2 []courseTime) bool {
	for _, t1 := range schedule1 {
		for _, t2 := range schedule2 {
			if isTimeConflict(t1, t2) {
				return true
			}
		}
	}
	return false
}

// isTimeConflict checks if two time slots conflict
func isTimeConflict(time1, time2 courseTime) bool {
	// Different days - no conflict
	if time1.Day != time2.Day {
		return false
	}

	// Same day - check time overlap
	// Convert to minutes for easier comparison
	start1 := time1.StartHour*60 + time1.StartMin
	end1 := time1.EndHour*60 + time1.EndMin
	start2 := time2.StartHour*60 + time2.StartMin
	end2 := time2.EndHour*60 + time2.EndMin

	// Conflict if: start1 < end2 AND start2 < end1
	return start1 < end2 && start2 < end1
}
