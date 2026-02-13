package cache

import (
	"course-reg/internal/app/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildConflictGraph(t *testing.T) {
	courses := []models.Course{
		{ID: 1, Schedules: "월 09:00~11:00", Capacity: 30},
		{ID: 2, Schedules: "월 10:00~12:00", Capacity: 30}, // conflicts with 1
		{ID: 3, Schedules: "수 14:00~16:00", Capacity: 30}, // no conflict
	}

	t.Run("conflict and no-conflict", func(t *testing.T) {
		graph, err := BuildConflictGraph(courses)
		assert.NoError(t, err)

		assert.True(t, graph[1][2])
		assert.True(t, graph[2][1])
		assert.False(t, graph[1][3])
		assert.False(t, graph[3][1])
		assert.False(t, graph[2][3])
		assert.False(t, graph[3][2])
	})

	t.Run("invalid schedule returns error", func(t *testing.T) {
		invalidCases := []struct {
			name     string
			schedule string
		}{
			{name: "empty string", schedule: ""},
			{name: "length too short", schedule: "월 9:10~11:30"},
			{name: "length too long", schedule: "월  09:10~11:30"},
			{name: "bad hour", schedule: "월 AB:10~11:30"},
			{name: "bad minute", schedule: "월 09:AB~11:30"},
		}

		for _, tt := range invalidCases {
			t.Run(tt.name, func(t *testing.T) {
				_, err := BuildConflictGraph([]models.Course{
					{ID: 1, Schedules: tt.schedule},
				})
				assert.Error(t, err)
			})
		}
	})
}

func TestIsTimeConflict(t *testing.T) {
	tests := []struct {
		name     string
		time1    courseTime
		time2    courseTime
		expected bool
	}{
		{
			name:     "different days - no conflict",
			time1:    courseTime{Day: "월", StartHour: 9, StartMin: 0, EndHour: 10, EndMin: 0},
			time2:    courseTime{Day: "화", StartHour: 9, StartMin: 0, EndHour: 10, EndMin: 0},
			expected: false,
		},
		{
			name:     "same day - full overlap",
			time1:    courseTime{Day: "월", StartHour: 9, StartMin: 0, EndHour: 10, EndMin: 0},
			time2:    courseTime{Day: "월", StartHour: 9, StartMin: 0, EndHour: 10, EndMin: 0},
			expected: true,
		},
		{
			name:     "same day - partial overlap",
			time1:    courseTime{Day: "월", StartHour: 9, StartMin: 0, EndHour: 10, EndMin: 30},
			time2:    courseTime{Day: "월", StartHour: 10, StartMin: 0, EndHour: 11, EndMin: 0},
			expected: true,
		},
		{
			name:     "same day - partial overlap (symmetric)",
			time1:    courseTime{Day: "월", StartHour: 10, StartMin: 0, EndHour: 11, EndMin: 0},
			time2:    courseTime{Day: "월", StartHour: 9, StartMin: 0, EndHour: 10, EndMin: 30},
			expected: true,
		},
		{
			name:     "same day - consecutive",
			time1:    courseTime{Day: "월", StartHour: 9, StartMin: 0, EndHour: 10, EndMin: 0},
			time2:    courseTime{Day: "월", StartHour: 10, StartMin: 0, EndHour: 11, EndMin: 0},
			expected: false,
		},
		{
			name:     "same day - consecutive (symmetric)",
			time1:    courseTime{Day: "월", StartHour: 10, StartMin: 0, EndHour: 11, EndMin: 0},
			time2:    courseTime{Day: "월", StartHour: 9, StartMin: 0, EndHour: 10, EndMin: 0},
			expected: false,
		},
		{
			name:     "same day - separated",
			time1:    courseTime{Day: "월", StartHour: 9, StartMin: 0, EndHour: 10, EndMin: 0},
			time2:    courseTime{Day: "월", StartHour: 14, StartMin: 0, EndHour: 15, EndMin: 0},
			expected: false,
		},
		{
			name:     "same day - separated (symmetric)",
			time1:    courseTime{Day: "월", StartHour: 14, StartMin: 0, EndHour: 15, EndMin: 0},
			time2:    courseTime{Day: "월", StartHour: 9, StartMin: 0, EndHour: 10, EndMin: 0},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isTimeConflict(tt.time1, tt.time2)
			if result != tt.expected {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}
