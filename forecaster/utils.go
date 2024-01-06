package forecaster

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"time"
)

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

func WriteSequencesCSV(writer io.Writer, sequences map[time.Time]uint) error {
	w := csv.NewWriter(writer)
	type Record struct {
		Date      time.Time `json:"date"`
		Occupancy uint      `json:"occupancy"`
	}
	records := make([]Record, 0, len(sequences))

	for date, occupiedSpots := range sequences {
		// record := []string{
		// 	date.Format(time.DateTime),
		// 	fmt.Sprintf("%f", float32(occupiedSpots)/float32(parkingLot.TotalSpots)),
		// }
		records = append(records, Record{
			Date:      date,
			Occupancy: occupiedSpots,
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
			fmt.Sprintf("%d", record.Occupancy),
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
