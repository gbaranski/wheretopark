package main

import (
	"net/url"
	"wheretopark/forecaster"
	wheretopark "wheretopark/go"
	"wheretopark/go/timeseries"

	"github.com/rs/zerolog/log"

	"github.com/caarlos0/env/v10"
)

type environment struct {
	Port           uint     `env:"PORT" envDefault:"8080"`
	TimeseriesData string   `env:"TIMESERIES_DATA,expand" envDefault:"${HOME}/.local/share/wheretopark/timeseries"`
	PycasterURL    *url.URL `env:"PYCASTER_URL,required"`
}

func main() {
	wheretopark.InitLogging()

	environment := environment{}
	if err := env.Parse(&environment); err != nil {
		log.Fatal().Err(err).Send()
	}

	timeseries := timeseries.New()
	err := timeseries.LoadMultipleCSV(environment.TimeseriesData)
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
