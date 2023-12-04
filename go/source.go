package wheretopark

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

type Source interface {
	ParkingLots(context.Context) (<-chan map[ID]ParkingLot, error)
}

func poll(ctx context.Context, cache *Cache, source Source) (uint, error) {
	ch, err := source.ParkingLots(ctx)
	if err != nil {
		return 0, err
	}
	count := uint(0)
	for parkingLots := range ch {
		count += uint(len(parkingLots))
		for id, parkingLot := range parkingLots {
			cache.SetParkingLot(id, &parkingLot)
		}
	}
	return count, nil
}

func RunPrefetch(cache *Cache, source Source, interval time.Duration) {
	ctx := log.With().Logger().WithContext(context.TODO())
	for {
		count, err := poll(ctx, cache, source)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("prefetch failed")
		} else {
			log.Ctx(ctx).Debug().Uint("parking-lots", count).Msg("prefetch succeesfull")
		}
		time.Sleep(interval)
	}
}
