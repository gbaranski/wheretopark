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
	"wheretopark/providers/collector/warsaw"

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

func runProvider(name string, createFn func() (wheretopark.Provider, error), client *wheretopark.Client) {
	provider, err := createFn()
	if err != nil {
		log.Fatalf("creating %s provider failed with error: %v", name, err)
	}
	err = wheretopark.RunProvider(client, provider)
	if err != nil {
		log.Fatalf("running %s provider failed with error: %v", name, err)
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

	go runProvider("gdansk", gdansk.NewProvider, client)
	go runProvider("gdynia", gdynia.NewProvider, client)
	go runProvider("warsaw", warsaw.NewProvider, client)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
