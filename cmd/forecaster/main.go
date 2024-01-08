package main

import (
	"net/url"
	"path/filepath"
	"wheretopark/forecaster"
	wheretopark "wheretopark/go"
	"wheretopark/go/timeseries"

	"github.com/rs/zerolog/log"

	"github.com/caarlos0/env/v10"
)

type environment struct {
	Port           uint     `env:"PORT" envDefault:"8080"`
	ForecasterData string   `env:"FORECASTER_DATA,expand" envDefault:"${HOME}/.local/share/wheretopark/forecaster"`
	PycasterURL    *url.URL `env:"PYCASTER_URL,required"`
}

func main() {
	wheretopark.InitLogging()

	environment := environment{}
	if err := env.Parse(&environment); err != nil {
		log.Fatal().Err(err).Send()
	}

	timeseries := timeseries.New()
	err := timeseries.LoadMultipleCSV(filepath.Join(environment.ForecasterData, "timeseries"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load timeseries")
	}

	pycaster := forecaster.NewPycaster(environment.PycasterURL)
	server := forecaster.NewServer(pycaster, timeseries)
	router := server.Router()
	if err := server.Run(router, environment.Port); err != nil {
		log.Fatal().Err(err).Msg("run server failure")
	}
}
