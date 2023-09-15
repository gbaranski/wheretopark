package main

// type config struct {
// 	DatabaseURL      string `env:"DATABASE_URL" envDefault:"ws://localhost:8000"`
// 	DatabaseName     string `env:"DATABASE_NAME" envDefault:"development"`
// 	DatabaseUser     string `env:"DATABASE_USER" envDefault:"root"`
// 	DatabasePassword string `env:"DATABASE_PASSWORD" envDefault:"password"`
// 	InfluxURL        string `env:"INFLUXDB_URL" envDefault:"http://localhost:8086"`
// 	InfluxToken      string `env:"INFLUXDB_TOKEN"`
// }

// func process(influx api.WriteAPIBlocking) error {
// 	parkingLots, err := client.GetAllParkingLots()
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("failed to fetch all parking lots")
// 	}

// 	for id, parkingLot := range parkingLots {
// 		point := influxdb2.NewPointWithMeasurement("parking_lot").
// 			AddTag("id", id).
// 			AddTag("name", parkingLot.Metadata.Name).
// 			AddField("availableSpots", parkingLot.State.AvailableSpots["CAR"]).
// 			AddField("totalSpots", parkingLot.Metadata.TotalSpots["CAR"]).
// 			SetTime(parkingLot.State.LastUpdated)

// 		err := influx.WritePoint(context.Background(), point)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	err = influx.Flush(context.Background())
// 	if err != nil {
// 		return fmt.Errorf("failed to flush: %w", err)
// 	}

// 	log.Info().Msgf("processed %d parking lots", len(parkingLots))
// 	return nil
// }

// func main() {
// 	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
// 	config := config{}
// 	if err := env.Parse(&config); err != nil {
// 		log.Fatal().Err(err).Send()
// 	}

// 	url, err := url.Parse(config.DatabaseURL)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("invalid database URL")
// 	}
// 	client, err := wheretopark.NewClient(url, "wheretopark", config.DatabaseName)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("failed to create database client")
// 	}
// 	defer client.Close()
// 	err = client.SignInWithPassword(config.DatabaseUser, config.DatabasePassword)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("failed to sign in")
// 	}

// 	influxClient := influxdb2.NewClient(config.InfluxURL, config.InfluxToken)
// 	defer influxClient.Close()
// 	writeAPI := influxClient.WriteAPIBlocking("wheretopark", config.DatabaseName)
// 	writeAPI.EnableBatching()

// 	for {
// 		err := process(client, writeAPI)
// 		if err != nil {
// 			log.Err(err).Msg("failed to process")
// 		}
// 		time.Sleep(time.Minute)
// 	}

// }
