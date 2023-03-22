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

func process(client *wheretopark.Client, provider Provider) error {
	metadatas, err := provider.GetMetadata()
	if err != nil {
		return fmt.Errorf("failed to get metadatas: %w", err)
	}
	log.Debug().
		Int("n", len(metadatas)).
		Str("name", provider.Name()).
		Msg("obtained metadatas")

	states, err := provider.GetState()
	if err != nil {
		return fmt.Errorf("failed to get states: %w", err)
	}
	log.Debug().
		Int("n", len(metadatas)).
		Str("name", provider.Name()).
		Msg("obtained states")

	for id, metadata := range metadatas {
		if metadata.PaymentMethods == nil {
			metadata.PaymentMethods = make([]string, 0)
		}
		if metadata.Comment == nil {
			metadata.Comment = make(map[string]string, 0)
		}
		for i, rule := range metadata.Rules {
			if rule.Pricing == nil {
				metadata.Rules[i].Pricing = make([]wheretopark.PricingRule, 0)
			}
		}
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

type Config struct {
	metadataInterval time.Duration
	stateInterval    time.Duration
}

var DEFAULT_CONFIG = Config{
	metadataInterval: time.Minute,
	stateInterval:    time.Minute,
}

func NewConfig(metadataInterval, stateInterval time.Duration) Config {
	return Config{
		metadataInterval: metadataInterval,
		stateInterval:    stateInterval,
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
		time.Sleep(config.metadataInterval)
	}
}
