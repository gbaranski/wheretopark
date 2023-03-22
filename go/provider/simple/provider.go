package simple

import (
	"fmt"
	"time"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
)

type Provider interface {
	Name() string
	Config() Config
	GetParkingLots() (map[wheretopark.ID]wheretopark.ParkingLot, error)
}

func process(client *wheretopark.Client, provider Provider) error {
	parkingLots, err := provider.GetParkingLots()
	if err != nil {
		return fmt.Errorf("failed to get data: %w", err)
	}
	log.Debug().
		Int("n", len(parkingLots)).
		Str("name", provider.Name()).
		Msg("obtained parking lots")

	err = client.SetParkingLots(parkingLots)
	if err != nil {
		return err
	}
	log.Info().
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
	log.Info().Str("name", provider.Name()).Msg("starting provider")
	config := provider.Config()
	for {
		processFn := func() error {
			return process(client, provider)
		}
		if err := wheretopark.WithTimeout(processFn, DEFAULT_PROCESS_TIMEOUT); err != nil {
			log.Error().
				Err(err).
				Str("name", provider.Name()).
				Msg("failed to process provider")
		}
		time.Sleep(config.interval)
	}
}
