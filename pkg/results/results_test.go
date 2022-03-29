package results

import (
	"testing"
	"time"
)

func TestCalculateMedianTime(t *testing.T) {
	l := []time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second, 4 * time.Second, 5 * time.Second}
	result := calculateMedianTime(l)
	if result != 3*time.Second {
		t.Errorf("calculateMedianTime was incorrect, got: %v, want: %v.", result, 3*time.Second)
	}
}

func TestCalculateTotalTime(t *testing.T) {
	l := []time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second, 4 * time.Second, 5 * time.Second}
	result := calculateTotalTime(l)
	if result != 15*time.Second {
		t.Errorf("calculateTotalTime was incorrect, got: %v, want: %v.", result, 3*time.Second)
	}
}
