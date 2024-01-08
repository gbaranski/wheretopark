package cctv

import wheretopark "wheretopark/go"

type Environment struct {
	ModelPath string           `env:"MODEL_PATH,expand" envDefault:"${HOME}/.local/share/wheretopark/model.onnx"`
	SavePath  *string          `env:"SAVE_PATH" envExpand:"true"`
	SaveItems []SaveItem       `env:"SAVE_ITEMS" envSeparator:","`
	SaveIDs   []wheretopark.ID `env:"SAVE_IDS" envSeparator:","`
}
