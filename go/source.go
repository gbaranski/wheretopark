package wheretopark

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

type Source interface {
	ParkingLots(context.Context) (<-chan map[ID]ParkingLot, error)
}

func RunPrefetch(cache *Cache, source Source, interval time.Duration) {
	for {
		ctx := log.With().Logger().WithContext(context.TODO())
		ch, err := source.ParkingLots(ctx)
		if err != nil {
			log.Error().Err(err).Msg("get parking lots failure")
			return
		}
		for parkingLots := range ch {
			for id, parkingLot := range parkingLots {
				cache.SetParkingLot(id, &parkingLot)
			}
		}
		log.Ctx(ctx).Debug().Msg("prefetch done")
		time.Sleep(interval)
	}
}
