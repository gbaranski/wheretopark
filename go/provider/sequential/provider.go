package sequential

import (
	"fmt"
	"time"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
)

type Provider interface {
	Name() string
	Config() Config
	GetMetadata() (map[wheretopark.ID]wheretopark.Metadata, error)
	GetState() (map[wheretopark.ID]wheretopark.State, error)
}

func obtainMetadatas(provider Provider, client *wheretopark.Client) (map[wheretopark.ID]wheretopark.Metadata, error) {
	metadatas, err := provider.GetMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to get metadatas: %w", err)
	}
	log.Debug().
		Int("n", len(metadatas)).
		Str("name", provider.Name()).
		Msg("obtained metadatas")
	return metadatas, nil
}

func obtainStates(provider Provider, client *wheretopark.Client) (map[wheretopark.ID]wheretopark.State, error) {
	states, err := provider.GetState()
	if err != nil {
		return nil, fmt.Errorf("failed to get states: %w", err)
	}
	log.Debug().
		Int("n", len(states)).
		Str("name", provider.Name()).
		Msg("obtained states")
	return states, nil
}

func processMetadata(provider Provider, client *wheretopark.Client) error {
	metadatas, err := obtainMetadatas(provider, client)
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

func processState(provider Provider, client *wheretopark.Client) error {
	states, err := obtainStates(provider, client)
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

func runLoop(processFn func() error, what string, interval time.Duration) {
	for {
		time.Sleep(interval)
		if err := wheretopark.WithTimeout(processFn, DEFAULT_PROCESS_TIMEOUT); err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("failed to process %s", what))
		}
	}
}

func Run(provider Provider, client *wheretopark.Client) error {
	log.Info().Str("name", provider.Name()).Msg("starting provider")
	config := provider.Config()

	err := initProcessing(provider, client)
	if err != nil {
		return fmt.Errorf("failed to initialize processing: %w", err)
	}

	go runLoop(func() error {
		return processState(provider, client)
	}, "state", config.stateInterval)
	runLoop(func() error {
		return processMetadata(provider, client)
	}, "metadata", config.metadataInterval)

	return nil
}

func initProcessing(provider Provider, client *wheretopark.Client) error {
	metadatas, err := obtainMetadatas(provider, client)
	if err != nil {
		return fmt.Errorf("failed to obtain metadatas: %w", err)
	}
	states, err := obtainStates(provider, client)
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
			return err
		}
	}
	log.Info().
		Int("n", len(metadatas)).
		Msg("updated parking lots")
	return nil
}
