package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	wheretopark "wheretopark/go"
	"wheretopark/providers/collector/gdansk"
	"wheretopark/providers/collector/gdynia"

	"github.com/caarlos0/env/v6"
)

type config struct {
	DatabaseURL         string  `env:"DATABASE_URL" envDefault:"ws://localhost:8000"`
	DatabaseName        string  `env:"DATABASE_NAME" envDefault:"development"`
	DatabaseUser        string  `env:"DATABASE_USER" envDefault:"root"`
	DatabasePassword    string  `env:"DATABASE_PASSWORD" envDefault:"root"`
	GdanskConfiguration *string `env:"GDANSK_CONFIGURATION"`
	GdyniaConfiguration *string `env:"GDYNIA_CONFIGURATION"`
}

func runGdansk(config *string, client *wheretopark.Client) {
	provider, err := gdansk.NewProvider(config)
	if err != nil {
		log.Fatalf("creating gdansk provider failed with error: %v", err)
	}
	err = wheretopark.RunProvider(client, provider)
	if err != nil {
		log.Fatalf("running gdansk provider failed with error: %v", err)
	}
}

func runGdynia(config *string, client *wheretopark.Client) {
	provider, err := gdynia.NewProvider(config)
	if err != nil {
		log.Fatalf("creating gdynia provider failed with error: %v", err)
	}
	err = wheretopark.RunProvider(client, provider)
	if err != nil {
		log.Fatalf("running gdynia provider failed with error: %v", err)
	}
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

	go runGdansk(config.GdanskConfiguration, client)
	go runGdynia(config.GdyniaConfiguration, client)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
