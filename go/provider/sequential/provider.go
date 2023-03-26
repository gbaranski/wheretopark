package sequential

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
	GetMetadata() (map[wheretopark.ID]wheretopark.Metadata, error)
	GetState() (map[wheretopark.ID]wheretopark.State, error)
}

func obtainMetadatas(logger zerolog.Logger, provider Provider, client *wheretopark.Client) (map[wheretopark.ID]wheretopark.Metadata, error) {
	metadatas, err := provider.GetMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to get metadatas: %w", err)
	}
	logger.Debug().
		Int("n", len(metadatas)).
		Msg("obtained metadatas")
	return metadatas, nil
}

func obtainStates(logger zerolog.Logger, provider Provider, client *wheretopark.Client) (map[wheretopark.ID]wheretopark.State, error) {
	states, err := provider.GetState()
	if err != nil {
		return nil, fmt.Errorf("failed to get states: %w", err)
	}
	logger.Debug().
		Int("n", len(states)).
		Msg("obtained states")
	return states, nil
}

func processMetadata(logger zerolog.Logger, provider Provider, client *wheretopark.Client) error {
	metadatas, err := obtainMetadatas(logger, provider, client)
	if err != nil {
		return err
	}

	for id, metadata := range metadatas {
		err := client.SetMetadata(id, metadata)
		if err != nil {
			return err
		}
	}
	return nil
}

func processState(logger zerolog.Logger, provider Provider, client *wheretopark.Client) error {
	states, err := obtainStates(logger, provider, client)
	if err != nil {
		return err
	}
	for id, state := range states {
		err := client.SetState(id, state)
		if err != nil {
			return err
		}
	}
	return nil
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

func runLoop(logger zerolog.Logger, processFn func() error, what string, interval time.Duration) {
	for {
		time.Sleep(interval)
		if err := wheretopark.WithTimeout(processFn, DEFAULT_PROCESS_TIMEOUT); err != nil {
			logger.Error().Err(err).Msg(fmt.Sprintf("failed to process %s", what))
		}
	}
}

func Run(provider Provider, client *wheretopark.Client) error {
	logger := log.With().Str("provider", provider.Name()).Logger()
	logger.Info().Str("type", "sequential").Msg("starting")
	config := provider.Config()

	err := initProcessing(logger, provider, client)
	if err != nil {
		return fmt.Errorf("failed to initialize processing: %w", err)
	}

	go runLoop(logger, func() error {
		return processState(logger, provider, client)
	}, "state", config.stateInterval)
	runLoop(logger, func() error {
		return processMetadata(logger, provider, client)
	}, "metadata", config.metadataInterval)
	return nil
}

func initProcessing(logger zerolog.Logger, provider Provider, client *wheretopark.Client) error {
	metadatas, err := obtainMetadatas(logger, provider, client)
	if err != nil {
		return fmt.Errorf("failed to obtain metadatas: %w", err)
	}
	states, err := obtainStates(logger, provider, client)
	if err != nil {
		return fmt.Errorf("failed to obtain metadatas: %w", err)
	}
	for id, metadata := range metadatas {
		parkingLot := wheretopark.ParkingLot{
			Metadata: metadata,
			State:    states[id],
		}
		err := client.SetParkingLot(id, parkingLot)
		if err != nil {
			return fmt.Errorf("setting parking lot: %w", err)
		}
	}
	logger.Info().
		Int("n", len(metadatas)).
		Msg("updated parking lots")
	return nil
}
