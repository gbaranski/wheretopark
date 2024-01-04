package forecaster

import (
	"fmt"
	"os"
	"sync"
	"time"
	wheretopark "wheretopark/go"

	"github.com/owulveryck/onnx-go"
	"github.com/owulveryck/onnx-go/backend/x/gorgonnx"
	"github.com/rs/zerolog/log"
	"gorgonia.org/tensor"
)

type Model struct {
	model   *onnx.Model
	backend *gorgonnx.Graph
	mu      *sync.Mutex
}

func NewModel(path string) (Model, error) {
	if _, err := os.Stat(path); err != nil {
		log.Fatal().Err(err).Msg("invalid model path")
	}
	backend := gorgonnx.NewGraph()
	model := onnx.NewModel(backend)
	b, err := os.ReadFile(path)
	if err != nil {
		return Model{}, fmt.Errorf("unable to read model %w", err)
	}
	err = model.UnmarshalBinary(b)
	if err != nil {
		return Model{}, fmt.Errorf("unable to unmarshal model %w", err)
	}

	return Model{
		model,
		backend,
		&sync.Mutex{},
	}, nil
}

const (
	NUM_STEPS    = 4
	NUM_FEATURES = 1
)

func (m *Model) Predict(parkingLot ParkingLot) (map[time.Time]uint, error) {
	sequences := SortedSequences(parkingLot.Sequences)
	segments := len(sequences) / NUM_STEPS
	input := make([]float32, segments*NUM_STEPS*NUM_FEATURES)
	for i, occupiedSpots := range sequences[:segments*NUM_STEPS] {
		rate := float32(occupiedSpots) / float32(parkingLot.TotalSpots)
		input[i] = rate
	}
	tensor := tensor.New(tensor.WithShape(segments, NUM_STEPS, NUM_FEATURES), tensor.WithBacking(input))
	m.mu.Lock()
	if err := m.model.SetInput(0, tensor); err != nil {
		return nil, fmt.Errorf("unable to set input %w", err)
	}
	if err := m.backend.Run(); err != nil {
		return nil, fmt.Errorf("unable to run model %w", err)
	}
	output, err := m.model.GetOutputTensors()
	if err != nil {
		return nil, fmt.Errorf("unable to get output %w", err)
	}
	for _, prediction := range output[0].Strides() {
		fmt.Printf("prediction: %d\n", prediction)
	}
	m.mu.Unlock()

	return nil, nil
}
