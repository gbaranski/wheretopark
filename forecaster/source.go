package forecaster

import (
	"time"
)

type Source interface {
	Name() string
	Load() (map[string]ParkingLot, error)
}

const DefaultInterval time.Duration = time.Minute * 15
const MinimumRecords uint = 50
