package wheretopark

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

type Source interface {
	ParkingLots(context.Context) (map[ID]ParkingLot, error)
}

type SequentialSource interface {
	Metadata(context.Context) (map[ID]Metadata, error)
	State(context.Context) (map[ID]State, error)
}

type SequentialSourceProxy struct {
	SequentialSource
}

func NewSequentialSourceProxy(source SequentialSource) SequentialSourceProxy {
	return SequentialSourceProxy{
		SequentialSource: source,
	}
}

func (s SequentialSourceProxy) ParkingLots(ctx context.Context) (map[ID]ParkingLot, error) {
	metadata, err := s.Metadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("get metadata fail: %w", err)
	}

	state, err := s.State(ctx)
	if err != nil {
		return nil, fmt.Errorf("get state fail: %w", err)
	}

	parkingLots := make(map[ID]ParkingLot)
	for id, metadata := range metadata {
		state, exists := state[id]
		if !exists {
			log.Debug().Str("id", id).Msg("missing state")
			continue
		}
		parkingLot := ParkingLot{
			State:    state,
			Metadata: metadata,
		}
		parkingLots[id] = parkingLot
	}
	return parkingLots, nil
}
