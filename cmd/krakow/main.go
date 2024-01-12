package main

import (
	wheretopark "wheretopark/go"
	"wheretopark/go/providers"
	"wheretopark/go/timeseries"
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

	timeseries := timeseries.New()
	err = timeseries.LoadMultipleCSV(environment.TimeseriesData)
	if err != nil {
		log.Fatal().Err(err).Msg("error loading timeseries data")
	}

	// sort.Slice(placemarks, func(i, j int) bool {
	// 	return placemarks[i].Code < placemarks[j].Code
	// })
	// for _, placemark := range placemarks {
	// 	googleMapsURL := fmt.Sprintf("https://www.google.com/maps/place/%f,%f", placemark.Coordinates.Latitude, placemark.Coordinates.Longitude)
	// 	fmt.Printf("%d: %s\n", placemark.Code, googleMapsURL)
	// }
	// return

	supportedPlacemarks := make([]krakow.Placemark, 0, len(placemarks))
	for _, placemark := range placemarks {
		id := placemark.ID()
		if timeseries.Contains(id) {
			supportedPlacemarks = append(supportedPlacemarks, placemark)
		}
	}

	provider := krakow.New(supportedPlacemarks, timeseries)

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
