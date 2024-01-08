package providers

import (
	"context"
	"time"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
)

type Provider interface {
	ParkingLots(context.Context) (<-chan map[wheretopark.ID]wheretopark.ParkingLot, error)
}

func poll(ctx context.Context, cache *Cache, provider Provider) (uint, error) {
	ch, err := provider.ParkingLots(ctx)
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

func RunPrefetch(cache *Cache, provider Provider, interval time.Duration) {
	ctx := log.With().Logger().WithContext(context.TODO())
	for {
		count, err := poll(ctx, cache, provider)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("prefetch failed")
		} else {
			log.Ctx(ctx).Debug().Uint("parking-lots", count).Msg("prefetch succeesfull")
		}
		time.Sleep(interval)
	}
}
