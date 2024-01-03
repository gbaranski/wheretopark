package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"wheretopark/forecaster"
	"wheretopark/forecaster/krakow"
	wheretopark "wheretopark/go"

	"github.com/caarlos0/env/v10"
	"github.com/rs/zerolog/log"
)

type environment struct {
	Source string `env:"FORECASTER_SOURCE,required"`
}

func main() {
	wheretopark.InitLogging()

	environment := environment{}
	if err := env.Parse(&environment); err != nil {
		log.Fatal().Err(err).Send()
	}
	datasetOutput := flag.String("output", "", "path to output the timeseries dataset")
	flag.Parse()

	if datasetOutput == nil || *datasetOutput == "" {
		log.Fatal().Msg("missing output path. specify with --output <path>")
	}

	krk := krakow.NewKrakow(filepath.Join(environment.Source, "krakow"))
	parkingLots, err := krk.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("error loading parking lots from krakow")
	}

	timeseries := forecaster.Timeseries{
		ParkingLots: parkingLots,
	}

	err = timeseries.SaveMultipleCSV(*datasetOutput)
	if err != nil {
		log.Fatal().Err(err).Msg("error saving output")
	}

	log.Info().Msg(fmt.Sprintf("wrote %d parking lots to %s", len(parkingLots), *datasetOutput))
}
