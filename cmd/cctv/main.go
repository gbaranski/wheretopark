package main

import (
	wheretopark "wheretopark/go"
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

	source := cctv.New(environment.CCTV)

	cache, err := wheretopark.NewCache()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create cache")
	}

	go wheretopark.RunPrefetch(cache, source, wheretopark.CacheTTL)

	server := wheretopark.NewServer(cache, source)
	router := server.Router()
	router.Get("/visualize/{id}/{camera}", source.HandleVisualizeCamera)
	if err := server.Run(router, environment.Port); err != nil {
		log.Fatal().Err(err).Msg("run server failure")
	}
}
