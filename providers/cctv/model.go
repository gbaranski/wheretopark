package cctv

import (
	"image"
	"sync"

	"gocv.io/x/gocv"
)

type Model struct {
	net gocv.Net
	mu  *sync.Mutex
}

func NewModel(path string) *Model {
	net := gocv.ReadNet(path, "")
	return &Model{
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
	predictions := make([]float32, len(images))
	for i, img := range images {
		prediction := m.Predict(img)
		predictions[i] = prediction
	}
	return predictions
}

func (m *Model) Close() {
	m.net.Close()
}
