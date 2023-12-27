package main

import (
	"context"
	"fmt"
	"time"
	wheretopark "wheretopark/go"

	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/caarlos0/env/v9"
	"github.com/rs/zerolog/log"
)

type environment struct {
	SERVER_URL   string `env:"SERVER_URL" envDefault:"https://api.wheretopark.app"`
	InfluxURL    string `env:"INFLUXDB_URL" envDefault:"https://eu-central-1-1.aws.cloud2.influxdata.com"`
	InfluxToken  string `env:"INFLUXDB_TOKEN"`
	InfluxBucket string `env:"INFLUXDB_BUCKET" envDefault:"parking_lots"`
}

func main() {
	wheretopark.InitLogging()

	environment := environment{}
	if err := env.Parse(&environment); err != nil {
		log.Fatal().Err(err).Send()
	}

	client := wheretopark.NewServerClient(wheretopark.MustParseURL(environment.SERVER_URL))

	// Create a new influx using an InfluxDB server base URL and an authentication token
	influx, err := influxdb3.New(influxdb3.ClientConfig{
		Host:  environment.InfluxURL,
		Token: environment.InfluxToken,
	})

	if err != nil {
		panic(err)
	}
	// Close client at the end and escalate error if present
	defer func(client *influxdb3.Client) {
		err := client.Close()
		if err != nil {
			panic(err)
		}
	}(influx)

	for {
		err := process(client, influx, environment.InfluxBucket)
		if err != nil {
			log.Err(err).Msg("failed to process")
		}
		time.Sleep(time.Minute * 15)
	}

}

func process(client *wheretopark.ServerClient, influx *influxdb3.Client, bucket string) error {
	providers, err := client.Providers()
	if err != nil {
		return fmt.Errorf("failed to fetch providers: %w", err)
	}

	parkingLots, err := client.GetFromMany(providers)
	if err != nil {
		return fmt.Errorf("failed to fetch all parking lots: %w", err)
	}

	points := make([]*influxdb3.Point, 0, len(parkingLots))
	for id, parkingLot := range parkingLots {
		point := influxdb3.NewPointWithMeasurement("availability").
			SetTag("id", id).
			SetTag("name", parkingLot.Metadata.Name).
			SetField("availableSpots", parkingLot.State.AvailableSpots["CAR"]).
			SetField("totalSpots", parkingLot.Metadata.TotalSpots["CAR"]).
			SetTimestamp(parkingLot.State.LastUpdated)
		points = append(points, point)
	}

	if err := influx.WritePointsWithOptions(context.Background(), &influxdb3.WriteOptions{
		Database: bucket,
	}, points...); err != nil {
		return fmt.Errorf("failed to write points: %w", err)
	}

	log.Info().Msgf("processed %d parking lots", len(points))
	return nil
}
