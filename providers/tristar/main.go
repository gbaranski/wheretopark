package main

import (
	"log"
	"net/url"
	wheretopark "wheretopark/go"
	"wheretopark/providers/tristar/gdansk"

	"github.com/caarlos0/env/v6"
)

type config struct {
	DatabaseURL         string  `env:"DATABASE_URL" envDefault:"ws://localhost:8000"`
	DatabaseName        string  `env:"DATABASE_NAME" envDefault:"development"`
	DatabaseUser        string  `env:"DATABASE_USER" envDefault:"root"`
	DatabasePassword    string  `env:"DATABASE_PASSWORD" envDefault:"root"`
	GdanskConfiguration *string `env:"GDANSK_CONFIGURATION"`
}

func main() {
	config := config{}
	if err := env.Parse(&config); err != nil {
		log.Fatalf("%+v\n", err)
	}

	url, err := url.Parse(config.DatabaseURL)
	if err != nil {
		log.Fatalf("invalid database url: %s", err)
	}
	client, err := wheretopark.NewClient(url, "wheretopark", config.DatabaseName)
	if err != nil {
		log.Fatalf("failed to create database client: %v", err)
	}
	defer client.Close()
	err = client.SignInWithPassword(config.DatabaseUser, config.DatabasePassword)
	if err != nil {
		log.Fatalf("failed to sign in: %v", err)
	}

	gdanskProvider, err := gdansk.NewProvider(config.GdanskConfiguration)
	if err != nil {
		log.Fatalf("creating provider failed with error: %v", err)
	}
	err = wheretopark.RunProvider(client, gdanskProvider)
	if err != nil {
		log.Fatalf("running provider failed with error: %v", err)

	}
}
