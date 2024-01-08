package meters

import (
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
	timeseries := timeseries.New()
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
	timeseries.FillMissingIntervals(opts.Interval)
	timeseries.Filter(func(id wheretopark.ID, interval time.Time) bool {
		return opts.WorkingScheme.Matches(interval)
	})
	timeseries.FilterTop(func(id wheretopark.ID, sequences map[time.Time]uint) bool {
		return len(sequences) >= int(opts.MinimumRecords)
	})

	return timeseries
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
