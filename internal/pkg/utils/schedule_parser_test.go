package utils

import (
	"testing"
)

func TestParseSchedule(t *testing.T) {
	t.Run("valid cases", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected []TimeSlot
		}{
			{
				name:  "single schedule",
				input: "월 09:10~11:30",
				expected: []TimeSlot{
					{Day: "월", StartHour: 9, StartMin: 10, EndHour: 11, EndMin: 30},
				},
			},
			{
				name:  "multiple schedules",
				input: "월 09:10~11:30, 수 17:10~19:20",
				expected: []TimeSlot{
					{Day: "월", StartHour: 9, StartMin: 10, EndHour: 11, EndMin: 30},
					{Day: "수", StartHour: 17, StartMin: 10, EndHour: 19, EndMin: 20},
				},
			},
			{
				name:  "all days",
				input: "월 09:00~10:00, 화 09:00~10:00, 수 09:00~10:00, 목 09:00~10:00, 금 09:00~10:00, 토 09:00~10:00, 일 09:00~10:00",
				expected: []TimeSlot{
					{Day: "월", StartHour: 9, StartMin: 0, EndHour: 10, EndMin: 0},
					{Day: "화", StartHour: 9, StartMin: 0, EndHour: 10, EndMin: 0},
					{Day: "수", StartHour: 9, StartMin: 0, EndHour: 10, EndMin: 0},
					{Day: "목", StartHour: 9, StartMin: 0, EndHour: 10, EndMin: 0},
					{Day: "금", StartHour: 9, StartMin: 0, EndHour: 10, EndMin: 0},
					{Day: "토", StartHour: 9, StartMin: 0, EndHour: 10, EndMin: 0},
					{Day: "일", StartHour: 9, StartMin: 0, EndHour: 10, EndMin: 0},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := ParseSchedule(tt.input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if len(result) != len(tt.expected) {
					t.Fatalf("length mismatch: got %d, want %d", len(result), len(tt.expected))
				}

				for i, slot := range result {
					exp := tt.expected[i]
					if slot.Day != exp.Day ||
						slot.StartHour != exp.StartHour ||
						slot.StartMin != exp.StartMin ||
						slot.EndHour != exp.EndHour ||
						slot.EndMin != exp.EndMin {
						t.Errorf("slot %d mismatch: got %+v, want %+v", i, slot, exp)
					}
				}
			})
		}
	})

	t.Run("error cases", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
		}{
			{
				name:  "empty string",
				input: "",
			},
			{
				name:  "invalid format - length too short",
				input: "월 9:10~11:30",
			},
			{
				name:  "invalid format - length too long",
				input: "월  09:10~11:30",
			},
			{
				name:  "invalid format - bad hour",
				input: "월 AB:10~11:30",
			},
			{
				name:  "invalid format - bad minute",
				input: "월 09:AB~11:30",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := ParseSchedule(tt.input)
				if err == nil {
					t.Errorf("expected error for input %q, got nil", tt.input)
				}
			})
		}
	})
}

func TestHasConflict(t *testing.T) {
	tests := []struct {
		name     string
		slot1    TimeSlot
		slot2    TimeSlot
		expected bool
	}{
		{
			name:     "different days - no conflict",
			slot1:    TimeSlot{Day: "월", StartHour: 9, StartMin: 0, EndHour: 10, EndMin: 0},
			slot2:    TimeSlot{Day: "화", StartHour: 9, StartMin: 0, EndHour: 10, EndMin: 0},
			expected: false,
		},
		{
			name:     "same day - full overlap",
			slot1:    TimeSlot{Day: "월", StartHour: 9, StartMin: 0, EndHour: 10, EndMin: 0},
			slot2:    TimeSlot{Day: "월", StartHour: 9, StartMin: 0, EndHour: 10, EndMin: 0},
			expected: true,
		},
		{
			name:     "same day - partial overlap",
			slot1:    TimeSlot{Day: "월", StartHour: 9, StartMin: 0, EndHour: 10, EndMin: 30},
			slot2:    TimeSlot{Day: "월", StartHour: 10, StartMin: 0, EndHour: 11, EndMin: 0},
			expected: true,
		},
		{
			name:     "same day - consecutive (no overlap)",
			slot1:    TimeSlot{Day: "월", StartHour: 9, StartMin: 0, EndHour: 10, EndMin: 0},
			slot2:    TimeSlot{Day: "월", StartHour: 10, StartMin: 0, EndHour: 11, EndMin: 0},
			expected: false,
		},
		{
			name:     "same day - separated",
			slot1:    TimeSlot{Day: "월", StartHour: 9, StartMin: 0, EndHour: 10, EndMin: 0},
			slot2:    TimeSlot{Day: "월", StartHour: 14, StartMin: 0, EndHour: 15, EndMin: 0},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasConflict(tt.slot1, tt.slot2)
			if result != tt.expected {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSchedulesConflict(t *testing.T) {
	t.Run("valid cases", func(t *testing.T) {
		tests := []struct {
			name      string
			schedule1 string
			schedule2 string
			expected  bool
		}{
			{
				name:      "conflict exists",
				schedule1: "월 09:00~10:30",
				schedule2: "월 10:00~11:00",
				expected:  true,
			},
			{
				name:      "no conflict - different days",
				schedule1: "월 09:00~10:30",
				schedule2: "화 09:00~10:30",
				expected:  false,
			},
			{
				name:      "no conflict - consecutive",
				schedule1: "월 09:00~10:00",
				schedule2: "월 10:00~11:00",
				expected:  false,
			},
			{
				name:      "multiple schedules - one conflicts",
				schedule1: "월 09:00~10:00, 수 14:00~15:00",
				schedule2: "화 09:00~10:00, 수 14:30~15:30",
				expected:  true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := SchedulesConflict(tt.schedule1, tt.schedule2)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("got %v, want %v", result, tt.expected)
				}
			})
		}
	})

	t.Run("error cases", func(t *testing.T) {
		tests := []struct {
			name      string
			schedule1 string
			schedule2 string
		}{
			{
				name:      "invalid schedule1",
				schedule1: "invalid",
				schedule2: "월 09:00~10:00",
			},
			{
				name:      "invalid schedule2",
				schedule1: "월 09:00~10:00",
				schedule2: "invalid",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := SchedulesConflict(tt.schedule1, tt.schedule2)
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			})
		}
	})
}
