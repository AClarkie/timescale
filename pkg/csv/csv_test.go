package csv

import (
	"testing"
)

func TestProcessCsv(t *testing.T) {
	_, err := ProcessCSV("a")
	if err == nil {
		t.Fatal("expected an error")
	}
}

// func TestGenerateHostnameMap(t *testing.T) {
// 	_, err := generateHostnameMap("a")
// 	if err == nil {
// 		t.Fatal("expected an error")
// 	}
// }
