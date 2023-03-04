package wheretopark

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

type Provider interface {
	Name() string
	GetMetadata() (map[ID]Metadata, error)
	GetState() (map[ID]State, error)
}

func process(client *Client, provider Provider) error {
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
				metadata.Rules[i].Pricing = make([]PricingRule, 0)
			}
		}
		parkingLot := ParkingLot{
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

type ProviderConfig struct {
	Interval time.Duration
}

var DEFAULT_PROVIDER_CONFIG = ProviderConfig{
	Interval: time.Minute,
}

func RunProvider(client *Client, provider Provider, config ProviderConfig) error {
	log.Info().Str("name", provider.Name()).Msg("starting provider")
	for {
		if err := process(client, provider); err != nil {
			log.Error().
				Err(err).
				Send()
		}
		time.Sleep(config.Interval)
	}
}
