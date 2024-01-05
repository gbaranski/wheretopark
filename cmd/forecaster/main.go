package main

import (
	"context"
	"fmt"
	"path/filepath"
	"wheretopark/forecaster"
	"wheretopark/forecaster/krakow"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"

	"github.com/caarlos0/env/v10"
)

// import (
// 	"flag"
// 	"fmt"
// 	"path/filepath"
// 	wheretopark "wheretopark/go"

// 	"github.com/caarlos0/env/v10"
// )

type environment struct {
	Data string `env:"FORECASTER_DATA,required"`
}

func main() {
	wheretopark.InitLogging()

	environment := environment{}
	if err := env.Parse(&environment); err != nil {
		log.Fatal().Err(err).Send()
	}

	fcSources := map[string]forecaster.ForecasterSource{
		"krakow": krakow.NewKrakow(filepath.Join(environment.Data, "datasets", "krakow")),
	}

	source := forecaster.New(fcSources)
	ch, err := source.ParkingLots(context.Background())
	if err != nil {
		panic(err)
	}
	for {
		parkingLot, ok := <-ch
		if !ok {
			break
		}
		fmt.Printf("parking lot: %+v\n", parkingLot)
	}
}
