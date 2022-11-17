package cctv

import (
	"image"

	"gocv.io/x/gocv"
)

type Model struct {
	net gocv.Net
}

func NewModel(path string) *Model {
	net := gocv.ReadNet(path, "")
	return &Model{
		net,
	}
}

func (m *Model) Predict(img gocv.Mat) float32 {
	blob := gocv.BlobFromImage(img, 1.0, image.Pt(128, 128), gocv.NewScalar(0, 0, 0, 0), true, false)
	defer blob.Close()
	blob.DivideUChar(255)

	m.net.SetInput(blob, "input_1")
	prob := m.net.Forward("")
	prediction := prob.GetFloatAt(0, 0)
	return prediction
}

func (m *Model) Close() {
	m.net.Close()
}
