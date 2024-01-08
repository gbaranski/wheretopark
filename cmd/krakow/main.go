package main

import (
	"fmt"
	"path/filepath"
	"strings"
	wheretopark "wheretopark/go"
	"wheretopark/go/providers"
	"wheretopark/providers/krakow"

	"github.com/caarlos0/env/v10"
	"github.com/rs/zerolog/log"
)

type environment struct {
	Port           uint   `env:"PORT" envDefault:"8080"`
	TimeseriesData string `env:"TIMESERIES_DATA,expand" envDefault:"${HOME}/.local/share/wheretopark/timeseries"`
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

	// this is probably just temporary, I have to filter out placemarks that are not in the timeseries data
	// so that user doesn't get them displayed on the map
	supportedPlacemarks := supportedPlacemarks(&environment, placemarks)
	fmt.Printf("supported placemarks: %d\n", len(supportedPlacemarks))

	provider := krakow.New(supportedPlacemarks)

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

func supportedPlacemarks(environment *environment, placemarks []krakow.Placemark) []krakow.Placemark {
	timeseriesFiles, err := wheretopark.ListFilesWithExtension(environment.TimeseriesData, "csv")
	if err != nil {
		log.Fatal().Err(err).Msg("error listing timeseries files")
	}
	supportedIDs := make(map[wheretopark.ID]struct{})
	for _, path := range timeseriesFiles {
		id := strings.TrimSuffix(filepath.Base(path), ".csv")
		supportedIDs[id] = struct{}{}
	}
	for id := range supportedIDs {
		fmt.Printf("supported placemark %s\n", id)
	}

	supportedPlacemarks := make([]krakow.Placemark, 0, len(placemarks))
	for _, placemark := range placemarks {
		id := placemark.ID()
		if _, ok := supportedIDs[id]; !ok {
			fmt.Printf("skipping placemark %s\n", id)
			continue
		}
		supportedPlacemarks = append(supportedPlacemarks, placemark)
	}
	return supportedPlacemarks
}
