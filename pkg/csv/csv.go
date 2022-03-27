package csv

import (
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

var (
	layout = "2006-01-02 15:04:05"
)

func ProcessCsv(queryParams string) (map[string]Hostname, error) {
	f, err := os.Open(queryParams)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	records := []*Record{}

	if err := gocsv.UnmarshalFile(f, &records); err != nil {
		return nil, err
	}

	hostnames := make(map[string]Hostname)

	for _, record := range records {
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
			TimeSlice := append(hostnames[record.Hostname].Times, timeRecord)
			hostnames[record.Hostname] = Hostname{Name: record.Hostname, Times: TimeSlice}
		} else {
			hostnames[record.Hostname] = Hostname{Name: record.Hostname, Times: []Time{timeRecord}}
		}
	}

	return hostnames, nil
}

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
