package csv

import (
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

var (
	layout = "2006-01-02 15:04:05"
)

type Record struct {
	Hostname  string `csv:"hostname"` // .csv column headers
	StartTime string `csv:"start_time"`
	EndTime   string `csv:"end_time"`
}

type Hostname struct {
	Name  string
	Times []Time
}

type Time struct {
	StartTime time.Time
	EndTime   time.Time
}

func ProcessCSV(queryParamsCSV string) (map[string]Hostname, error) {
	f, err := os.Open(queryParamsCSV)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	records, err := readCSV(f)
	if err != nil {
		return nil, err
	}

	return generateHostnameMap(records)
}

func readCSV(f *os.File) ([]*Record, error) {
	records := []*Record{}
	if err := gocsv.UnmarshalFile(f, &records); err != nil {
		return nil, err
	}
	return records, nil
}

// Group the time entries by hostname into a map of hostname (key) to Hostname type (value)
func generateHostnameMap(records []*Record) (map[string]Hostname, error) {
	hostnames := make(map[string]Hostname)

	for _, record := range records {
		// Parse the time fields
		startTime, err := time.Parse(layout, record.StartTime)
		if err != nil {
			return nil, err
		}
		endTime, err := time.Parse(layout, record.EndTime)
		if err != nil {
			return nil, err
		}

		timeRecord := Time{StartTime: startTime, EndTime: endTime}

		if _, ok := hostnames[record.Hostname]; ok {
			Times := append(hostnames[record.Hostname].Times, timeRecord)
			hostnames[record.Hostname] = Hostname{Name: record.Hostname, Times: Times}
		} else {
			hostnames[record.Hostname] = Hostname{Name: record.Hostname, Times: []Time{timeRecord}}
		}
	}
	return hostnames, nil
}
