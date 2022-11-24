package wheretopark

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

type Provider interface {
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
		Msg("obtained metadatas")

	states, err := provider.GetState()
	if err != nil {
		return fmt.Errorf("failed to get states: %w", err)
	}
	log.Debug().
		Int("n", len(metadatas)).
		Msg("obtained states")

	for id, metadata := range metadatas {
		parkingLot := ParkingLot{
			Metadata: metadata,
			State:    states[id],
		}
		err := client.SetParkingLot(id, parkingLot)
		if err != nil {
			return err
		}
		log.Debug().Str("id", id).Msg("updated parking lot")
	}
	log.Info().
		Int("n", len(metadatas)).
		Msg("updated parking lots")

	return nil
}

func RunProvider(client *Client, provider Provider) error {
	for {
		if err := process(client, provider); err != nil {
			log.Error().
				Err(err).
				Send()
		}
		time.Sleep(time.Minute)
	}
}
