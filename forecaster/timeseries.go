package forecaster

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
)

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

	w := csv.NewWriter(f)
	defer w.Flush()
	records := make([][]string, 0, len(parkingLot.Sequences)+1)
	records = append(records, []string{
		"date",
		"occupancy",
	})

	type StructSequence struct {
		Date          time.Time
		OccupiedSpots uint
	}

	sequences := make([]StructSequence, 0, len(parkingLot.Sequences))
	for date, occupiedSpots := range parkingLot.Sequences {
		sequences = append(sequences, StructSequence{
			Date:          date,
			OccupiedSpots: occupiedSpots,
		})
	}
	sort.Slice(sequences, func(i, j int) bool {
		return sequences[i].Date.Before(sequences[j].Date)
	})
	for _, seq := range sequences {
		record := []string{
			seq.Date.Format(time.DateTime),
			fmt.Sprintf("%f", float32(seq.OccupiedSpots)/float32(parkingLot.TotalSpots)),
		}
		records = append(records, record)
	}

	for _, record := range records {
		if err := w.Write(record); err != nil {
			return fmt.Errorf("error writing record to file: %w", err)
		}
	}
	return nil
}
