package main

import (
	wheretopark "wheretopark/go"
	"wheretopark/go/providers"
	"wheretopark/providers/krakow"

	"github.com/caarlos0/env/v10"
	"github.com/rs/zerolog/log"
)

type environment struct {
	Port uint `env:"PORT" envDefault:"8080"`
}

func main() {
	wheretopark.InitLogging()

	environment := environment{}
	if err := env.Parse(&environment); err != nil {
		log.Fatal().Err(err).Send()
	}

	placemarks, err := krakow.GetPlacemarks()
	if err != nil {
		log.Fatal().Err(err).Msg("error getting placemarks")
	}

	provider := krakow.New(placemarks)

	cache, err := providers.NewCache()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create cache")
	}

	// go wheretopark.RunPrefetch(cache, provider, wheretopark.CacheTTL)

	server := providers.NewServer(cache, provider)
	router := server.Router()
	if err := server.Run(router, environment.Port); err != nil {
		log.Fatal().Err(err).Msg("run server failure")
	}
}
