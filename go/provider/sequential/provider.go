package sequential

import (
	"time"
	wheretopark "wheretopark/go"
)

type Provider interface {
	Name() string
	Config() Config
	GetMetadatas() (map[wheretopark.ID]wheretopark.Metadata, error)
	GetStates() (map[wheretopark.ID]wheretopark.State, error)
}

type Config struct {
	metadataInterval time.Duration
	stateInterval    time.Duration
}

var DEFAULT_CONFIG = Config{
	metadataInterval: time.Minute * 5,
	stateInterval:    time.Minute,
}

func NewConfig(metadataInterval, stateInterval time.Duration) Config {
	return Config{
		metadataInterval: metadataInterval,
		stateInterval:    stateInterval,
	}
}

const DEFAULT_PROCESS_TIMEOUT = 30 * time.Second
