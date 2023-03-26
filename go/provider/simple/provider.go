package simple

import (
	"fmt"
	"time"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Provider interface {
	Name() string
	Config() Config
	GetParkingLots() (map[wheretopark.ID]wheretopark.ParkingLot, error)
}

func process(logger zerolog.Logger, client *wheretopark.Client, provider Provider) error {
	parkingLots, err := provider.GetParkingLots()
	if err != nil {
		return fmt.Errorf("failed to get data: %w", err)
	}
	logger.Debug().
		Int("n", len(parkingLots)).
		Msg("obtained parking lots")

	err = client.SetParkingLots(parkingLots)
	if err != nil {
		return err
	}
	logger.Info().
		Int("n", len(parkingLots)).
		Msg("updated parking lots")

	return nil
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

func Run(provider Provider, client *wheretopark.Client) error {
	logger := log.With().Str("provider", provider.Name()).Logger()
	logger.Info().Str("type", "simple").Msg("starting")
	config := provider.Config()
	for {
		processFn := func() error {
			return process(logger, client, provider)
		}
		if err := wheretopark.WithTimeout(processFn, DEFAULT_PROCESS_TIMEOUT); err != nil {
			logger.Error().
				Err(err).
				Msg("failed to process provider")
		}
		time.Sleep(config.interval)
	}
}
