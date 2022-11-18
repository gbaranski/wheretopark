package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	wheretopark "wheretopark/go"
	"wheretopark/providers/cctv"

	"github.com/caarlos0/env/v6"
)

type environment struct {
	DatabaseURL      string  `env:"DATABASE_URL" envDefault:"ws://localhost:8000"`
	DatabaseName     string  `env:"DATABASE_NAME" envDefault:"development"`
	DatabaseUser     string  `env:"DATABASE_USER" envDefault:"root"`
	DatabasePassword string  `env:"DATABASE_PASSWORD" envDefault:"root"`
	Configuration    *string `env:"CONFIGURATION"`
	Model            string  `env:"MODEL" envDefault:"$HOME/.local/share/wheretopark/cctv/model.onnx" envExpand:"true"`
}

func main() {
	environment := environment{}
	if err := env.Parse(&environment); err != nil {
		log.Fatalf("%+v\n", err)
	}

	url, err := url.Parse(environment.DatabaseURL)
	if err != nil {
		log.Fatalf("invalid database url: %s", err)
	}
	client, err := wheretopark.NewClient(url, "wheretopark", environment.DatabaseName)
	if err != nil {
		log.Fatalf("failed to create database client: %v", err)
	}
	defer client.Close()
	err = client.SignInWithPassword(environment.DatabaseUser, environment.DatabasePassword)
	if err != nil {
		log.Fatalf("failed to sign in: %v", err)
	}

	model := cctv.NewModel(environment.Model)
	defer model.Close()
	provider, err := cctv.NewProvider(environment.Configuration, model)
	if err != nil {
		panic(err)
	}
	defer provider.Close()
	err = wheretopark.RunProvider(client, provider)
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
