package main

import (
	"log"
	wheretopark "wheretopark/go"

	"github.com/caarlos0/env/v6"
)

type config struct {
	DatabaseURL      string `env:"DATABASE_URL"`
	DatabaseName     string `env:"DATABASE_NAME"`
	DatabaseUser     string `env:"DATABASE_USER"`
	DatabasePassword string `env:"DATABASE_PASSWORD"`
}

func main() {
	config := config{}
	if err := env.Parse(&config); err != nil {
		log.Fatalf("%+v\n", err)
	}

	client, err := wheretopark.NewClient(config.DatabaseURL, "wheretopark", "development")
	if err != nil {
		log.Fatalf("failed to create database client: %v", err)
	}
	err = client.SignInWithPassword(config.DatabaseUser, config.DatabasePassword)
	if err != nil {
		log.Fatalf("failed to sign in: %v", err)
	}

	// gdyniaProvider := gdansk.Provider{}

}
