package forecaster

import (
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
