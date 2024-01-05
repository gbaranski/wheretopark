package forecaster

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
)

type ParkingLot struct {
	wheretopark.Metadata
	Sequences map[time.Time]uint
}

type SequenceTime = time.Time

func MaxOccupiedSpots(sequences map[time.Time]uint) uint {
	count := uint(0)
	for _, value := range sequences {
		count = max(count, uint(value))
	}
	return count
}

func SortedSequences(sequencesMap map[time.Time]uint) []uint {
	// Create a slice of keys (dates) from the map
	dates := make([]time.Time, 0, len(sequencesMap))
	for date := range sequencesMap {
		dates = append(dates, date)
	}
	// Sort the slice of dates
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	// Create a slice of sorted OccupiedSpots
	sortedSequences := make([]uint, 0, len(dates))
	for _, date := range dates {
		sortedSequences = append(sortedSequences, sequencesMap[date])
	}

	return sortedSequences
}

func (p ParkingLot) WriteCSV(writer io.Writer) error {
	w := csv.NewWriter(writer)
	type Record struct {
		Date      time.Time `json:"date"`
		Occupancy float32   `json:"occupancy"`
	}
	records := make([]Record, 0, len(p.Sequences))

	for date, occupiedSpots := range p.Sequences {
		// record := []string{
		// 	date.Format(time.DateTime),
		// 	fmt.Sprintf("%f", float32(occupiedSpots)/float32(parkingLot.TotalSpots)),
		// }
		records = append(records, Record{
			Date:      date,
			Occupancy: float32(occupiedSpots) / float32(p.AllTotalSpots()),
		})
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].Date.Before(records[j].Date)
	})

	w.Write([]string{
		"date",
		"occupancy",
	})
	for _, record := range records {
		recordString := []string{
			record.Date.Format(time.DateTime),
			fmt.Sprintf("%f", record.Occupancy),
		}
		if err := w.Write(recordString); err != nil {
			return fmt.Errorf("error writing record to file: %w", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return fmt.Errorf("error flushing writer: %w", err)
	}
	return nil
}

type Timeseries struct {
	ParkingLots map[wheretopark.ID]ParkingLot `json:"parkingLots"`
}

func (t Timeseries) SaveJSON(basePath string) error {
	path := filepath.Join(basePath, "timeseries.json")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer f.Close()

	jsonData, err := json.MarshalIndent(t, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshalling timeseries: %w", err)
	}
	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing json: %w", err)
	}

	return nil
}

func (t Timeseries) SaveMultipleCSV(basePath string) error {
	for id, parkingLot := range t.ParkingLots {
		path := filepath.Join(basePath, fmt.Sprintf("%s.csv", id))
		err := SaveSingleCSV(path, parkingLot)
		if err != nil {
			log.Error().Err(err).Str("id", id).Msg("error saving csv")
		}
	}

	return nil
}

func SaveSingleCSV(path string, parkingLot ParkingLot) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer f.Close()

	if err = parkingLot.WriteCSV(f); err != nil {
		return fmt.Errorf("error writing csv: %w", err)
	}
	return nil
}
