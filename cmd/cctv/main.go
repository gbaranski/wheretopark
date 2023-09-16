package main

import (
	wheretopark "wheretopark/go"
	"wheretopark/providers/cctv"

	"github.com/caarlos0/env/v8"
	"github.com/rs/zerolog/log"
)

type environment struct {
	Port      int              `env:"PORT" envDefault:"8080"`
	ModelPath string           `env:"MODEL_PATH" envDefault:"$HOME/.local/share/wheretopark/cctv/model.onnx" envExpand:"true"`
	SavePath  *string          `env:"SAVE_PATH" envExpand:"true"`
	SaveItems []cctv.SaveItem  `env:"SAVE_ITEMS" envSeparator:","`
	SaveIDs   []wheretopark.ID `env:"SAVE_IDS" envSeparator:","`
}

func main() {
	wheretopark.InitLogging()

	environment := environment{}
	if err := env.Parse(&environment); err != nil {
		log.Fatal().Err(err).Send()
	}

	provider, err := cctv.NewProvider(environment.ModelPath, environment.SavePath, environment.SaveItems, environment.SaveIDs)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer provider.Close()

	if err := wheretopark.RunProvider(provider, uint(environment.Port)); err != nil {
		log.Fatal().Err(err).Msg("run provider failure")
	}
}
