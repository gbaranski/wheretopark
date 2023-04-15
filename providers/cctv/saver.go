package cctv

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
	wheretopark "wheretopark/go"

	"gocv.io/x/gocv"
)

type Saver struct {
	basePath *string
	items    []SaveItem
	ids      []wheretopark.ID
}

func NewSaver(basePath *string, items []SaveItem, ids []wheretopark.ID) Saver {
	return Saver{
		basePath: basePath,
		items:    items,
		ids:      ids,
	}
}

func (s *Saver) SavePredictions(img gocv.Mat, id string, cameraID int, captureTime time.Time, spots []ParkingSpot, predictions []float32) error {
	if s.basePath == nil || (len(s.ids) > 0 && !contains(s.ids, id)) {
		return nil
	}
	basePath := fmt.Sprintf("%s/%s/%02d", *s.basePath, id, cameraID)
	time := captureTime.UTC().Format("2006-01-02--15-04-05")
	savers := map[SaveItem]func() error{
		SaveItemInput: func() error {
			return saveInput(img, fmt.Sprintf("%s/inputs/%s.jpg", basePath, time))
		},
		SaveItemResult: func() error {
			return saveResults(spots, predictions, fmt.Sprintf("%s/results/%s.json", basePath, time))
		},
		SaveItemVisualization: func() error {
			return saveVisualizations(img, spots, predictions, fmt.Sprintf("%s/visualizations/%s.jpg", basePath, time))
		},
	}
	for _, saver := range s.items {
		if err := savers[saver](); err != nil {
			return fmt.Errorf("failed saving %s: %w", saver, err)
		}
	}
	return nil
}

type SpotResult struct {
	Prediction float32 `json:"prediction"`
	Points     []Point `json:"points"`
}

type Result struct {
	Spots []SpotResult `json:"spots"`
}

type SaveItem = string

const (
	SaveItemVisualization SaveItem = "visualizations"
	SaveItemResult        SaveItem = "results"
	SaveItemInput         SaveItem = "inputs"
)

func saveInput(img gocv.Mat, path string) error {
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", directory, err)
	}
	ok := gocv.IMWrite(path, img)
	if !ok {
		return fmt.Errorf("failed to write input image to %s", path)
	}
	return nil
}

func saveResults(spots []ParkingSpot, predictions []float32, path string) error {
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", directory, err)
	}
	spotResults := make([]SpotResult, len(predictions))
	for i, prediction := range predictions {
		spot := spots[i]
		spotResults[i] = SpotResult{
			Prediction: prediction,
			Points:     spot.Points,
		}
	}
	result := Result{
		Spots: spotResults,
	}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}
	err = os.WriteFile(path, resultJSON, 0644)
	if err != nil {
		return fmt.Errorf("failed to write result to %s: %w", path, err)
	}
	return nil
}

func saveVisualizations(img gocv.Mat, spots []ParkingSpot, predictions []float32, path string) error {
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", directory, err)
	}
	for i, prediction := range predictions {
		spot := spots[i]
		VisualizeSpotPrediction(&img, spot, prediction)
	}
	ok := gocv.IMWrite(path, img)
	if !ok {
		return fmt.Errorf("failed to write visualization to %s", path)
	}
	return nil
}
