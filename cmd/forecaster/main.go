package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
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

	err = SaveOutput(parkingLots, *datasetOutput)
	if err != nil {
		log.Fatal().Err(err).Msg("error saving output")
	}
}

func SaveOutput(parkingLots map[wheretopark.ID]forecaster.ParkingLot, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	timeseries := forecaster.Timeseries{
		ParkingLots: parkingLots,
	}
	jsonData, err := json.MarshalIndent(timeseries, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshalling timeseries: %w", err)
	}
	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing json: %w", err)
	}
	log.Info().Msg(fmt.Sprintf("wrote %d parking lots to %s", len(timeseries.ParkingLots), path))

	return nil
}
