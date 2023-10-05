package main

import (
	"context"
	"time"
	"wheretopark/collector"
	"wheretopark/collector/cctv"
	"wheretopark/collector/gdansk"
	"wheretopark/collector/gdynia"
	"wheretopark/collector/glasgow"
	"wheretopark/collector/lacity"
	"wheretopark/collector/poznan"
	"wheretopark/collector/warsaw"
	wheretopark "wheretopark/go"

	"github.com/caarlos0/env/v9"
	"github.com/rs/zerolog/log"
)

type environment struct {
	Port uint             `env:"PORT" envDefault:"8080"`
	CCTV cctv.Environment `envPrefix:"CCTV_"`
}

// useful when source is slow
func prefetch(cache *wheretopark.Cache, sources map[string]wheretopark.Source, interval time.Duration) {
	for {
		for name, source := range sources {
			ctx := log.With().Str("source", name).Logger().WithContext(context.TODO())
			ch, err := source.ParkingLots(ctx)
			if err != nil {
				log.Error().Err(err).Str("name", name).Msg("get parking lots failure")
				continue
			}
			for parkingLots := range ch {
				cache.UpdateParkingLots(name, parkingLots)
			}
			log.Ctx(ctx).Debug().Msg("prefetch done")
		}
		time.Sleep(interval)
	}
}

func main() {
	wheretopark.InitLogging()

	environment := environment{}
	if err := env.Parse(&environment); err != nil {
		log.Fatal().Err(err).Send()
	}

	sources := map[string]wheretopark.Source{
		"gdansk":  gdansk.New(),
		"gdynia":  gdynia.New(),
		"glasgow": glasgow.New(),
		"lacity":  lacity.New(),
		"poznan":  poznan.New(),
		"warsaw":  warsaw.New(),
		"cctv":    cctv.New(environment.CCTV),
	}

	cache, err := wheretopark.NewCache()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create cache")
	}

	go prefetch(cache, map[string]wheretopark.Source{
		"cctv": sources["cctv"],
	}, wheretopark.CacheTTL)

	server := collector.NewServer(cache, sources)
	router := server.Router()
	router.Get("/visualize/{id}/{camera}", sources["cctv"].(cctv.Source).HandleVisualizeCamera)
	if err := server.Run(router, environment.Port); err != nil {
		log.Fatal().Err(err).Msg("run server failure")
	}
}
