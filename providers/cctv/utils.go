package cctv

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"time"

	"gocv.io/x/gocv"
)

func (spot *ParkingSpot) ImagePoints() []image.Point {
	points := make([]image.Point, len(spot.Points))
	for i, point := range spot.Points {
		points[i] = point.ToImagePoint()
	}
	return points
}

// PS: frame is image
func CropSpot(cvFrame gocv.Mat, spot *ParkingSpot) gocv.Mat {
	// obtain transformation
	points := spot.ImagePoints()
	cvPoints := gocv.NewPointVectorFromPoints(points)
	defer cvPoints.Close()
	cvMinAreaRect := gocv.MinAreaRect(cvPoints)
	size := image.Point{X: cvMinAreaRect.Width, Y: cvMinAreaRect.Height}
	cvBoxMat := gocv.NewMat()
	defer cvBoxMat.Close()
	gocv.BoxPoints(cvMinAreaRect, &cvBoxMat)
	box := []image.Point{}
	for i := 0; i < cvBoxMat.Rows(); i++ {
		x := cvBoxMat.GetFloatAt(i, 0)
		y := cvBoxMat.GetFloatAt(i, 1)
		point := image.Pt(int(x), int(y))
		box = append(box, point)
	}

	cvBox := gocv.NewPointVectorFromPoints(box)
	defer cvBox.Close()
	cvDestination := gocv.NewPointVectorFromPoints(
		[]image.Point{
			{X: 0, Y: size.Y - 1},
			{X: 0, Y: 0},
			{X: size.X - 1, Y: 0},
			{X: size.X - 1, Y: size.Y - 1},
		},
	)
	defer cvDestination.Close()
	cvTransformation := gocv.GetPerspectiveTransform(cvBox, cvDestination)
	defer cvTransformation.Close()

	// transform
	cvOutput := gocv.NewMat()
	gocv.WarpPerspective(cvFrame, &cvOutput, cvTransformation, size)
	return cvOutput
}

func VisualizeSpotPrediction(img *gocv.Mat, spot ParkingSpot, prediction float32) {
	occupied := prediction < 0.5
	var drawingColor color.RGBA
	if occupied {
		drawingColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	} else {
		drawingColor = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	}
	points := spot.ImagePoints()
	cvPoints := gocv.NewPointVectorFromPoints(points)
	defer cvPoints.Close()
	cvPointsVector := gocv.NewPointsVectorFromPoints([][]image.Point{points})
	defer cvPointsVector.Close()
	cvMinAreaRect := gocv.MinAreaRect(cvPoints)
	gocv.Polylines(img, cvPointsVector, true, drawingColor, 2)
	gocv.PutText(img, fmt.Sprintf("%.2f", prediction), cvMinAreaRect.Center, gocv.FontHersheyPlain, 1.0, drawingColor, 1)
}

func ExtractSpots(img gocv.Mat, spots []ParkingSpot) []gocv.Mat {
	images := make([]gocv.Mat, len(spots))
	for i, spot := range spots {
		images[i] = CropSpot(img, &spot)
	}
	return images
}

func SaveRawImage(img gocv.Mat, path string) error {
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", directory, err)
	}
	ok := gocv.IMWrite(path, img)
	if !ok {
		return fmt.Errorf("failed to write raw image to %s", path)
	}
	return nil
}

func SaveResults(spots []ParkingSpot, predictions []float32, path string) error {
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

func SaveVisualizations(img gocv.Mat, spots []ParkingSpot, predictions []float32, path string) error {
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

func SavePredictions(img gocv.Mat, basePath string, captureTime time.Time, spots []ParkingSpot, predictions []float32) error {
	time := captureTime.UTC().Format("2006-01-02--15-04-05")
	if err := SaveRawImage(img, fmt.Sprintf("%s/images/%s.jpg", basePath, time)); err != nil {
		return err
	}
	if err := SaveResults(spots, predictions, fmt.Sprintf("%s/results/%s.json", basePath, time)); err != nil {
		return err
	}
	if err := SaveVisualizations(img, spots, predictions, fmt.Sprintf("%s/visualizations/%s.jpg", basePath, time)); err != nil {
		return err
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
