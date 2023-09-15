package simple

import (
	"time"
	wheretopark "wheretopark/go"
)

type Provider interface {
	Name() string
	Config() Config
	GetParkingLots() (map[wheretopark.ID]wheretopark.ParkingLot, error)
}

type Config struct {
	interval time.Duration
}

var DEFAULT_CONFIG = Config{
	interval: time.Minute,
}

func NewConfig(interval time.Duration) Config {
	return Config{
		interval: interval,
	}
}

const DEFAULT_PROCESS_TIMEOUT = 30 * time.Second
