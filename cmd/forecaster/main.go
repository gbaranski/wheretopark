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
	krkMetadata, err := krk.Metadata()
	if err != nil {
		log.Fatal().Err(err).Msg("error loading metadata")
	}
	meters, err := krk.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("error loading meters")
	}

	metersByID := make(map[wheretopark.ID]*forecaster.ParkingMeter)
	for code, meter := range meters {
		id := ""
		for _, metadata := range krkMetadata {
			if metadata.Code == code {
				if _, ok := metersByID[code]; ok {
					log.Error().Str("code", code).Msg("duplicate meter")
					break
				}
				id = wheretopark.CoordinateToID(metadata.Coordinates.Latitude, metadata.Coordinates.Longitude)
			}
		}
		if id == "" {
			log.Warn().Str("code", code).Msg("meter without metadata")
			continue
		}
		metersByID[id] = meter
	}

	// for id, meter := range metersByID {
	// 	log.Info().Str("id", id).Str("name", meter.Name).Int("occupancyData", len(meter.OccupancyData)).Uint("totalSpots", meter.TotalSpots()).Msg("meter")
	// }
	err = SaveOutput(metersByID, *datasetOutput)
	if err != nil {
		log.Fatal().Err(err).Msg("error saving output")
	}
}

func SaveOutput(meters map[wheretopark.ID]*forecaster.ParkingMeter, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	timeseries := forecaster.NewTimeseries(meters)
	jsonData, err := json.Marshal(timeseries)
	if err != nil {
		return fmt.Errorf("error marshalling timeseries: %w", err)
	}
	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing json: %w", err)
	}
	log.Info().Msg(fmt.Sprintf("wrote %d parking lots to %s", len(timeseries.ParkingLots), path))

	// records := timeseries.EncodeCSV()
	// for _, record := range records {
	// 	if err := w.Write(record); err != nil {
	// 		return fmt.Errorf("error writing record: %w", err)
	// 	}
	// }
	// log.Info().Msg(fmt.Sprintf("wrote %d records to %s", len(records), path))
	return nil
}
