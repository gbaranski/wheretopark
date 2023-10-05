package cctv

import (
	"image"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
	"gocv.io/x/gocv"
)

type Model struct {
	net gocv.Net
	mu  *sync.Mutex
}

func NewModel(path string) Model {
	if _, err := os.Stat(path); err != nil {
		log.Fatal().Err(err).Msg("invalid model path")
	}
	net := gocv.ReadNetFromONNX(path)
	return Model{
		net,
		&sync.Mutex{},
	}
}

func (m *Model) Predict(img gocv.Mat) float32 {
	blob := gocv.BlobFromImage(img, 1.0, image.Pt(128, 128), gocv.NewScalar(0, 0, 0, 0), true, false)
	defer blob.Close()
	blob.DivideUChar(255)

	m.mu.Lock()
	defer m.mu.Unlock()

	m.net.SetInput(blob, "input_1")
	prob := m.net.Forward("")
	defer prob.Close()
	prediction := prob.GetFloatAt(0, 0)
	return prediction

}

func (m *Model) PredictMany(images []gocv.Mat) []float32 {
	blob := gocv.NewMat()
	defer blob.Close()
	gocv.BlobFromImages(images, &blob, 1.0, image.Pt(128, 128), gocv.NewScalar(0, 0, 0, 0), true, false, gocv.MatTypeCV8UC3)
	blob.DivideUChar(255)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.net.SetInput(blob, "input_1")
	prob := m.net.Forward("")
	defer prob.Close()
	predictions := make([]float32, len(images))
	for i := 0; i < prob.Size()[0]; i++ {
		predictions[i] = prob.GetFloatAt(i, 0)
	}
	return predictions
}

func (m *Model) Close() {
	m.net.Close()
}
