package main

import (
	"net/url"
	"os"
	"os/signal"
	"syscall"
	wheretopark "wheretopark/go"
	"wheretopark/go/provider"
	"wheretopark/go/provider/sequential"
	"wheretopark/go/provider/simple"
	"wheretopark/providers/collector/gdansk"
	"wheretopark/providers/collector/gdynia"
	"wheretopark/providers/collector/glasgow"
	"wheretopark/providers/collector/poznan"
	"wheretopark/providers/collector/warsaw"

	"github.com/caarlos0/env/v7"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type config struct {
	DatabaseURL      string `env:"DATABASE_URL" envDefault:"ws://localhost:8000"`
	DatabaseName     string `env:"DATABASE_NAME" envDefault:"development"`
	DatabaseUser     string `env:"DATABASE_USER" envDefault:"root"`
	DatabasePassword string `env:"DATABASE_PASSWORD" envDefault:"password"`
}

func runProvider[T provider.Common](createFn func() (T, error), runFn func(T, *wheretopark.Client) error, client *wheretopark.Client) {
	provider, err := createFn()
	name := provider.Name()
	if err != nil {
		log.Fatal().Err(err).Str("provider", name).Msg("fail to create")
	}
	err = runFn(provider, client)
	if err != nil {
		log.Fatal().Err(err).Str("provider", name).Msg("fail to run")
	}
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	config := config{}
	if err := env.Parse(&config); err != nil {
		log.Fatal().Err(err).Send()
	}

	url, err := url.Parse(config.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid database URL")
	}
	client, err := wheretopark.NewClient(url, "wheretopark", config.DatabaseName)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create database client")
	}
	defer client.Close()
	err = client.SignInWithPassword(config.DatabaseUser, config.DatabasePassword)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to sign in")
	}

	go runProvider(gdansk.NewProvider, sequential.Run, client)
	go runProvider(gdynia.NewProvider, sequential.Run, client)
	go runProvider(warsaw.NewProvider, simple.Run, client)
	go runProvider(poznan.NewProvider, simple.Run, client)
	// TODO: After getting the business subscription, decrease the interval
	go runProvider(glasgow.NewProvider, simple.Run, client)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
