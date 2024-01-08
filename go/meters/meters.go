package meters

import (
	"os"
	"path/filepath"
	"strings"
	"time"
	wheretopark "wheretopark/go"
	"wheretopark/go/timeseries"
)

const DefaultMinimumRecords uint = 50
const DefaultInterval time.Duration = time.Minute * 15

type WorkingScheme struct {
	StartHour uint8
	EndHour   uint8
	Weekdays  []time.Weekday
}

type Options struct {
	MinimumRecords uint
	Interval       time.Duration
	WorkingScheme  WorkingScheme
	Mapping        map[uint]wheretopark.ID
}

type Record struct {
	Code      uint
	StartDate time.Time
	EndDate   time.Time
}

func NewTimeseriesFromMeterRecords(records []Record, opts Options) timeseries.TimeSeries {
	timeseries := timeseries.NewTimeseries(opts.Interval)
	for _, record := range records {
		id, exists := opts.Mapping[record.Code]
		if !exists {
			continue
		}

		currentInterval := record.StartDate.Truncate(opts.Interval)
		endInterval := record.EndDate.Truncate(opts.Interval)
		for currentInterval.Before(endInterval) {
			timeseries.Add(id, currentInterval)
			currentInterval = currentInterval.Add(opts.Interval)
		}
	}
	timeseries.FillMissingIntervals()
	timeseries.Filter(func(id wheretopark.ID, interval time.Time) bool {
		return opts.WorkingScheme.Matches(interval)
	})
	timeseries.FilterTop(func(id wheretopark.ID, sequences map[time.Time]uint) bool {
		return len(sequences) >= int(opts.MinimumRecords)
	})

	return timeseries
}

func ListFilesWithExtension(basePath string, extension string) ([]string, error) {
	var filesWithExtension []string

	// Read all files and directories within basePath
	items, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	// Iterate through the items
	for _, item := range items {
		if !item.IsDir() { // Ensure it's a file, not a directory
			// Check if the file has the desired extension
			if strings.HasSuffix(item.Name(), "."+extension) {
				fullPath := filepath.Join(basePath, item.Name())
				filesWithExtension = append(filesWithExtension, fullPath)
			}
		}
	}

	return filesWithExtension, nil
}

func (s WorkingScheme) Matches(t time.Time) bool {
	startHourAtDay := time.Date(t.Year(), t.Month(), t.Day(), int(s.StartHour), 0, 0, 0, time.UTC)
	endHourAtDay := time.Date(t.Year(), t.Month(), t.Day(), int(s.EndHour), 0, 0, 0, time.UTC)
	if t.Before(startHourAtDay) || t.After(endHourAtDay) {
		return false
	}

	for _, day := range s.Weekdays {
		if t.Weekday() == day {
			return true
		}
	}
	return false
}
