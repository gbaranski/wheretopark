package main

import (
	"net/http"
	"net/url"
	"os"
	wheretopark "wheretopark/go"
	"wheretopark/prometheus/exporter"

	"github.com/caarlos0/env/v7"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type config struct {
	DatabaseURL      string `env:"DATABASE_URL" envDefault:"ws://localhost:8000"`
	DatabaseName     string `env:"DATABASE_NAME" envDefault:"development"`
	DatabaseUser     string `env:"DATABASE_USER" envDefault:"root"`
	DatabasePassword string `env:"DATABASE_PASSWORD" envDefault:"password"`
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

	collector := exporter.NewCollector(*client)
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	log.Info().Msg("Starting exporter on :2112")
	http.ListenAndServe(":2112", nil)
}
