package main

import (
	wheretopark "wheretopark/go"

	"wheretopark/providers/collector"

	"github.com/caarlos0/env/v8"
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

	provider, err := collector.NewProvider()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create provider")
	}
	if err := wheretopark.RunProvider(provider, uint(environment.Port)); err != nil {
		log.Fatal().Err(err).Msg("run provider failure")
	}
}
