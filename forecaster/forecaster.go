package forecaster

import (
	"sort"
	"time"
)

type ParkingLot struct {
	TotalSpots uint               `json:"totalSpots"`
	Sequences  map[time.Time]uint `json:"sequences"`
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
