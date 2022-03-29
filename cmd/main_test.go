package main

import (
	"testing"
	"time"

	"github.com/AClarkie/timescale/pkg/csv"
)

func TestAssignHostsToGroups(t *testing.T) {
	pTime, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
	time := csv.Time{StartTime: pTime, EndTime: pTime}
	hostnames := map[string]csv.Hostname{
		"1": {Name: "1", Times: []csv.Time{time}},
		"2": {Name: "2", Times: []csv.Time{time}},
		"3": {Name: "3", Times: []csv.Time{time}},
		"4": {Name: "4", Times: []csv.Time{time}},
	}

	result := assignHostsToGroups(hostnames, 1)

	if len(result) != 1 {
		t.Errorf("assignHostsToGroups was incorrect, got: %d, want: %d.", len(result), 1)
	}

	result = assignHostsToGroups(hostnames, 4)

	if len(result) != 4 {
		t.Errorf("assignHostsToGroups was incorrect, got: %d, want: %d.", len(result), 4)
	}

	result = assignHostsToGroups(hostnames, 2)

	if len(result) != 2 {
		t.Errorf("assignHostsToGroups was incorrect, got: %d, want: %d.", len(result), 2)
	}
}
