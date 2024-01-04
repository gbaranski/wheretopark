package forecaster_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"wheretopark/forecaster"
	"wheretopark/forecaster/krakow"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
)

var (
	TESTED_PARKING_ID   = "u2yhvbnb0e"
	FORECASTER_PATH     = filepath.Join(wheretopark.Must(os.UserHomeDir()), ".local/share/wheretopark/forecaster")
	KRAKOW_DATASET_PATH = filepath.Join(FORECASTER_PATH, "datasets", "krakow")
	PARKING_MODEL_PATH  = filepath.Join(FORECASTER_PATH, "models", fmt.Sprintf("%s.onnx", TESTED_PARKING_ID))
)

func TestModel(t *testing.T) {
	model, err := forecaster.NewModel(PARKING_MODEL_PATH)
	if err != nil {
		t.Fatal(err)
	}
	krk := krakow.NewKrakow(KRAKOW_DATASET_PATH)
	parkingLots, err := krk.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("error loading parking lots from krakow")
	}
	parkingLot := parkingLots[TESTED_PARKING_ID]
	predictions, err := model.Predict(parkingLot)
	if err != nil {
		t.Fatal(err)
	}
	t.Fatal(predictions)
}
