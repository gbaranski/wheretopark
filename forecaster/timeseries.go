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
	type Record struct {
		Date      time.Time `json:"date"`
		Occupancy float32   `json:"occupancy"`
	}
	records := make([]Record, 0, len(parkingLot.Sequences))

	for date, occupiedSpots := range parkingLot.Sequences {
		// record := []string{
		// 	date.Format(time.DateTime),
		// 	fmt.Sprintf("%f", float32(occupiedSpots)/float32(parkingLot.TotalSpots)),
		// }
		records = append(records, Record{
			Date:      date,
			Occupancy: float32(occupiedSpots) / float32(parkingLot.TotalSpots),
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
	return nil
}
