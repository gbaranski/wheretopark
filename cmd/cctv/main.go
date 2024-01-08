package main

import (
	wheretopark "wheretopark/go"
	"wheretopark/go/providers"
	"wheretopark/providers/cctv"

	"github.com/caarlos0/env/v10"
	"github.com/rs/zerolog/log"
)

type environment struct {
	Port uint             `env:"PORT" envDefault:"8080"`
	CCTV cctv.Environment `envPrefix:"CCTV_"`
}

func main() {
	wheretopark.InitLogging()

	environment := environment{}
	if err := env.Parse(&environment); err != nil {
		log.Fatal().Err(err).Send()
	}

	provider := cctv.New(environment.CCTV)

	cache, err := providers.NewCache()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create cache")
	}

	go providers.RunPrefetch(cache, provider, providers.CacheTTL)

	server := providers.NewServer(cache, provider)
	router := server.Router()
	router.Get("/visualize/{id}/{camera}", provider.HandleVisualizeCamera)
	if err := server.Run(router, environment.Port); err != nil {
		log.Fatal().Err(err).Msg("run server failure")
	}
}
