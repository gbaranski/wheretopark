package forecaster

import (
	"time"
)

type Source interface {
	Name() string
	Load() (map[string]ParkingLot, error)
}

const DefaultInterval time.Duration = time.Hour
const MinimumRecords uint = 50
