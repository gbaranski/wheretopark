package main

import (
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
	wheretopark "wheretopark/go"
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
	DatabaseURL         string  `env:"DATABASE_URL" envDefault:"ws://localhost:8000"`
	DatabaseName        string  `env:"DATABASE_NAME" envDefault:"development"`
	DatabaseUser        string  `env:"DATABASE_USER" envDefault:"root"`
	DatabasePassword    string  `env:"DATABASE_PASSWORD" envDefault:"root"`
	GdanskConfiguration *string `env:"GDANSK_CONFIGURATION"`
	GdyniaConfiguration *string `env:"GDYNIA_CONFIGURATION"`
}

func runProvider(createFn func() (wheretopark.Provider, error), client *wheretopark.Client, config wheretopark.ProviderConfig) {
	provider, err := createFn()
	name := provider.Name()
	if err != nil {
		log.Fatal().Err(err).Str("name", name).Msg("creating provider failed")
	}
	err = wheretopark.RunProvider(client, provider, config)
	if err != nil {
		log.Fatal().Err(err).Str("name", name).Msg("running provider failed")
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

	go runProvider(gdansk.NewProvider, client, wheretopark.DEFAULT_PROVIDER_CONFIG)
	go runProvider(gdynia.NewProvider, client, wheretopark.DEFAULT_PROVIDER_CONFIG)
	go runProvider(warsaw.NewProvider, client, wheretopark.DEFAULT_PROVIDER_CONFIG)
	go runProvider(poznan.NewProvider, client, wheretopark.DEFAULT_PROVIDER_CONFIG)
  go runProvider(glasgow.NewProvider, client, wheretopark.ProviderConfig{Interval: 5 * time.Minute})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
