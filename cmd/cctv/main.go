package main

import (
	"fmt"
	wheretopark "wheretopark/go"
	"wheretopark/providers/cctv"

	"github.com/caarlos0/env/v8"
	"github.com/rs/zerolog/log"
	"gocv.io/x/gocv"
)

type environment struct {
	Port      int              `env:"PORT" envDefault:"8080"`
	ModelPath string           `env:"MODEL_PATH" envDefault:"$HOME/.local/share/wheretopark/cctv/yolov8x.onnx" envExpand:"true"`
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
	// saver := cctv.NewSaver(environment.SavePath, environment.SaveItems, environment.SaveIDs)

	img, err := generateVisualizationFor(model, &cctv.DefaultConfiguration.ParkingLots["u2gyfvc23d"].Cameras[0])
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer img.Close()
	window := gocv.NewWindow("Where To Park")
	window.IMShow(img)
	window.WaitKey(1)

	// provider, err := cctv.NewProvider(model, saver)
	// if err != nil {
	// 	log.Fatal().Err(err).Send()
	// }
	// defer provider.Close()

	// r, err := wheretopark.GetProviderRouter(provider)
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("get router failure")
	// }

	// for id, parkingLot := range cctv.DefaultConfiguration.ParkingLots {
	// 	log.Info().Str("id", id).Str("name", parkingLot.Name).Int("cameras", len(parkingLot.Cameras)).Msg("registering parking lot")
	// }

	// r.Get("/visualize/{id}/{camera}", func(w http.ResponseWriter, r *http.Request) {
	// 	id := wheretopark.ID(chi.URLParam(r, "id"))
	// 	parkingLot, exists := cctv.DefaultConfiguration.ParkingLots[id]
	// 	if !exists {
	// 		w.WriteHeader(http.StatusNotFound)
	// 		w.Write([]byte("parking lot not found"))
	// 		return
	// 	}

	// 	cameraID, err := strconv.Atoi(chi.URLParam(r, "camera"))
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		w.Write([]byte("invalid camera ID"))
	// 		return
	// 	}
	// 	if cameraID > len(parkingLot.Cameras)-1 {
	// 		w.WriteHeader(http.StatusNotFound)
	// 		w.Write([]byte("camera not found"))
	// 		return
	// 	}
	// 	camera := parkingLot.Cameras[cameraID]

	// 	img, err := generateVisualizationFor(model, &camera)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		w.Write([]byte(fmt.Sprintf("generate visualization failure: %s", err)))
	// 		return
	// 	}
	// 	defer img.Close()

	// 	buf, err := gocv.IMEncode(".jpg", img)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		log.Err(err).Msg("encode image failure")
	// 		return
	// 	}
	// 	defer buf.Close()

	// 	w.Header().Set("Content-Type", "image/jpeg")
	// 	w.Write(buf.GetBytes())
	// })

	// if err := wheretopark.RunRouter(r, uint(environment.Port)); err != nil {
	// 	log.Fatal().Err(err).Msg("run provider failure")
	// }
}
