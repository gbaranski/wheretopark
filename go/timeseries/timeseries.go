package timeseries

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
)

type TimeSeries struct {
	sequences map[wheretopark.ID]map[time.Time]uint
}

func New() TimeSeries {
	return TimeSeries{
		sequences: make(map[wheretopark.ID]map[time.Time]uint),
	}
}

func (ts *TimeSeries) CountFor(id wheretopark.ID) uint {
	return uint(len(ts.sequences[id]))
}

func (ts *TimeSeries) IDs() []wheretopark.ID {
	ids := make([]wheretopark.ID, 0, len(ts.sequences))
	for id := range ts.sequences {
		ids = append(ids, id)
	}
	return ids
}

func (ts *TimeSeries) Get(id wheretopark.ID) map[time.Time]uint {
	seq, ok := ts.sequences[id]
	if ok {
		return seq
	} else {
		return nil
	}
}

func (ts *TimeSeries) Add(id wheretopark.ID, interval time.Time) {
	if _, ok := ts.sequences[id]; !ok {
		ts.sequences[id] = make(map[time.Time]uint)
	}
	ts.sequences[id][interval]++
}

func (ts *TimeSeries) FillMissingIntervals(interval time.Duration) {
	for id, sequences := range ts.sequences {
		earliest := time.Now()
		latest := time.Time{}
		for interval := range sequences {
			if interval.Before(earliest) {
				earliest = interval
			}
			if interval.After(latest) {
				latest = interval
			}
		}
		currentInterval := earliest
		endInterval := latest
		for currentInterval.Before(endInterval) {
			if _, ok := ts.sequences[id][currentInterval]; !ok {
				ts.sequences[id][currentInterval] = 0
			}
			currentInterval = currentInterval.Add(interval)
		}
	}
}

func (ts *TimeSeries) MaxOccupancyOf(id wheretopark.ID) uint {
	maximum := 0
	for _, count := range ts.sequences[id] {
		maximum = max(maximum, int(count))
	}
	return uint(maximum)
}

func (ts *TimeSeries) FilterTop(predicate func(id wheretopark.ID, sequences map[time.Time]uint) bool) {
	for id, sequences := range ts.sequences {
		if !predicate(id, sequences) {
			delete(ts.sequences, id)
		}
	}
}

func (ts *TimeSeries) Filter(predicate func(id wheretopark.ID, interval time.Time) bool) {
	for id, sequences := range ts.sequences {
		for interval := range sequences {
			if !predicate(id, interval) {
				delete(ts.sequences[id], interval)
			}
		}
	}
}

func (ts *TimeSeries) writeSingleCSV(id wheretopark.ID, path string) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		log.Info().Str("path", path).Msg("file already exists, overwriting")
	}
	sequences := ts.sequences[id]
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	err = writer.Write([]string{"date", "occupancy"})
	if err != nil {
		return fmt.Errorf("failed to write headers to csv file: %w", err)
	}

	// Write data
	for t, occupancy := range sequences {
		err = writer.Write([]string{t.Format(time.RFC3339), strconv.FormatUint(uint64(occupancy), 10)})
		if err != nil {
			return fmt.Errorf("failed to write record to csv file: %w", err)
		}
	}
	return nil

}

func (ts *TimeSeries) WriteMultipleCSV(basePath string) error {
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		err = os.MkdirAll(basePath, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	for id := range ts.sequences {
		path := filepath.Join(basePath, fmt.Sprintf("%s.csv", id))
		err := ts.writeSingleCSV(id, path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ts *TimeSeries) loadSingleCSV(path string) (map[time.Time]uint, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	sequences := make(map[time.Time]uint)

	// Read and discard headers
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read headers from csv file: %w", err)
	}

	// Read data
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading csv file: %w", err)
		}

		t, err := time.Parse(time.RFC3339, record[0])
		if err != nil {
			return nil, fmt.Errorf("error parsing time: %w", err)
		}
		u, err := strconv.ParseUint(record[1], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("error parsing uint: %w", err)
		}

		sequences[t] = uint(u)
	}

	return sequences, nil

}

func (ts *TimeSeries) LoadMultipleCSV(basePath string) error {
	files, err := wheretopark.ListFilesWithExtension(basePath, "csv")
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}
	for _, path := range files {
		id := strings.TrimSuffix(path, ".csv")
		sequences, err := ts.loadSingleCSV(path)
		if err != nil {
			log.Error().Err(err).Str("path", path).Msg("failed to load csv file")
		}
		ts.sequences[id] = sequences
	}
	log.Info().Msg(fmt.Sprintf("loaded %d sequences from %d files", len(ts.sequences), len(files)))
	return nil
}
