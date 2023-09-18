package main

import (
	"fmt"
	"net/http"
	"strconv"
	wheretopark "wheretopark/go"
	"wheretopark/providers/cctv"

	"github.com/caarlos0/env/v8"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"gocv.io/x/gocv"
)

type environment struct {
	Port      int              `env:"PORT" envDefault:"8080"`
	ModelPath string           `env:"MODEL_PATH" envDefault:"$HOME/.local/share/wheretopark/cctv/model.onnx" envExpand:"true"`
	SavePath  *string          `env:"SAVE_PATH" envExpand:"true"`
	SaveItems []cctv.SaveItem  `env:"SAVE_ITEMS" envSeparator:","`
	SaveIDs   []wheretopark.ID `env:"SAVE_IDS" envSeparator:","`
}

func generateVisualizationFor(model cctv.Model, camera *cctv.ParkingLotCamera) (gocv.Mat, error) {
	img, err := cctv.GetImageFromCamera(camera.URL)
	if err != nil {
		return gocv.NewMat(), fmt.Errorf("unable to get image from camera: %v", err)
	}

	spotImages := cctv.ExtractSpots(img, camera.Spots)
	defer func() {
		for _, spotImage := range spotImages {
			spotImage.Close()
		}
	}()

	predictions := model.PredictMany(spotImages)
	for i, prediction := range predictions {
		spot := camera.Spots[i]
		cctv.VisualizeSpotPrediction(&img, spot, prediction)
	}
	return img, nil
}

func main() {
	wheretopark.InitLogging()

	environment := environment{}
	if err := env.Parse(&environment); err != nil {
		log.Fatal().Err(err).Send()
	}

	model := cctv.NewModel(environment.ModelPath)
	saver := cctv.NewSaver(environment.SavePath, environment.SaveItems, environment.SaveIDs)

	provider, err := cctv.NewProvider(model, saver)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer provider.Close()

	r, err := wheretopark.GetProviderRouter(provider)
	if err != nil {
		log.Fatal().Err(err).Msg("get router failure")
	}

	r.GET("/visualize/:id/:camera", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := wheretopark.ID(p.ByName("id"))
		parkingLot, exists := cctv.DefaultConfiguration.ParkingLots[id]
		if !exists {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("parking lot not found"))
			return
		}

		cameraID, err := strconv.Atoi(p.ByName("camera"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid camera ID"))
			return
		}
		var camera *cctv.ParkingLotCamera
		for i, currentCamera := range parkingLot.Cameras {
			if i+1 == cameraID {
				camera = &currentCamera
			}
		}
		if camera == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("camera not found"))
			return
		}

		img, err := generateVisualizationFor(model, camera)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("generate visualization failure: %s", err)))
			return
		}
		defer img.Close()

		buf, err := gocv.IMEncode(".jpg", img)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Err(err).Msg("encode image failure")
			return
		}
		defer buf.Close()

		w.Header().Set("Content-Type", "image/jpeg")
		w.Write(buf.GetBytes())
	})

	if err := wheretopark.RunRouter(r, uint(environment.Port)); err != nil {
		log.Fatal().Err(err).Msg("run provider failure")
	}
}
