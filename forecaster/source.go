package forecaster

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
)

const DefaultInterval time.Duration = time.Minute * 15
const MinimumRecords uint = 50

type ForecasterSource interface {
	Load() (map[wheretopark.ID]ParkingLot, error)
}

type CollectorSource struct {
	sources  map[string]ForecasterSource
	pycaster Pycaster
}

func (s CollectorSource) ParkingLots(ctx context.Context) (<-chan map[wheretopark.ID]wheretopark.ParkingLot, error) {
	ch := make(chan map[wheretopark.ID]wheretopark.ParkingLot, 1)
	var wg sync.WaitGroup
	for name, source := range s.sources {
		parkingLots, err := source.Load()
		if err != nil {
			log.Fatal().Err(err).Str("source", name).Msg("error loading parking lots from")
		}
		for id, parkingLot := range parkingLots {
			wg.Add(1)
			go func(id wheretopark.ID, parkingLot ParkingLot) {
				defer wg.Done()
				predictions, err := s.pycaster.Predict(id, parkingLot)
				if err != nil {
					log.Error().Err(err).Str("id", id).Msg("error predicting occupancy")
				}
				ch <- map[wheretopark.ID]wheretopark.ParkingLot{
					id: wheretopark.ParkingLot{
						Metadata: parkingLot.Metadata,
						State:    wheretopark.State{},
					},
				}
			}(id, parkingLot)
			return nil, nil
		}
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	return nil, nil
}

var (
	BASE_PATH = filepath.Join(wheretopark.Must(os.UserHomeDir()), ".local/share/wheretopark/forecaster")
)

func New(sources map[string]ForecasterSource, pycaster Pycaster) CollectorSource {
	return CollectorSource{
		sources:  sources,
		pycaster: pycaster,
	}
}
