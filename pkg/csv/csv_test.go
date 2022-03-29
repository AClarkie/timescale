package csv

import (
	"testing"
)

func TestProcessNonExistentCsv(t *testing.T) {
	_, err := ProcessCSV("a")
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestGenerateHostnameMap(t *testing.T) {
	records := []*Record{{Hostname: "1", StartTime: "2006-01-02 15:04:05", EndTime: "2006-01-02 15:04:05"},
		{Hostname: "1", StartTime: "2007-01-02 15:04:05", EndTime: "2007-01-02 15:04:05"}}

	result, _ := generateHostnameMap(records)

	hostname := result["1"]
	if len(hostname.Times) != 2 {
		t.Errorf("generateHostnameMap incorrect number of Time entries, got: %d, want: %d.", len(hostname.Times), 2)
	}
}
