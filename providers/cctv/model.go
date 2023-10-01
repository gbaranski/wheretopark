package cctv

import (
	"image"
	"sync"

	"github.com/rs/zerolog/log"
	"gocv.io/x/gocv"
)

type Model struct {
	net gocv.Net
	mu  *sync.Mutex
}

func NewModel(path string) Model {
	net := gocv.ReadNetFromONNX(path)
	log.Info().Str("path", path).Msg("loaded onnx model")
	return Model{
		net,
		&sync.Mutex{},
	}
}

const (
	HEIGHT = 640
	WIDTH  = 640
)

var SIZE = image.Pt(HEIGHT, WIDTH)

func (m *Model) Predict(img gocv.Mat) float32 {
	gocv.Resize(img, &img, SIZE, 0, 0, gocv.InterpolationLanczos4)

	blob := gocv.BlobFromImage(img, 1.0, SIZE, gocv.NewScalar(0, 0, 0, 0), true, false)
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
