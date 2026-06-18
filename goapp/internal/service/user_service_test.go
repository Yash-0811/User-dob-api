package service_test

import (
	"testing"
	"time"

	"github.com/yash/user-dob-api/internal/service"
)

func TestCalculateAge(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		dob      time.Time
		expected int
	}{
		{
			name:     "birthday already passed this year",
			dob:      now.AddDate(-30, -1, 0), // 30 years ago, 1 month ago
			expected: 30,
		},
		{
			name:     "birthday is tomorrow",
			dob:      now.AddDate(-25, 0, 1), // 25 years ago, birthday tomorrow
			expected: 24,
		},
		{
			name:     "birthday is today",
			dob:      now.AddDate(-20, 0, 0),
			expected: 20,
		},
		{
			name:     "newborn (< 1 year)",
			dob:      now.AddDate(0, -6, 0),
			expected: 0,
		},
		{
			name:     "exactly 1 year ago",
			dob:      now.AddDate(-1, 0, 0),
			expected: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := service.CalculateAge(tc.dob)
			if got != tc.expected {
				t.Errorf("CalculateAge(%v) = %d, want %d", tc.dob, got, tc.expected)
			}
		})
	}
}
