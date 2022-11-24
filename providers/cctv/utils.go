package cctv

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"os"

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
	cvMinAreaRect := gocv.MinAreaRect(cvPoints)
	size := image.Point{X: cvMinAreaRect.Width, Y: cvMinAreaRect.Height}
	cvBoxMat := gocv.NewMat()
	gocv.BoxPoints(cvMinAreaRect, &cvBoxMat)
	box := []image.Point{}
	for i := 0; i < cvBoxMat.Rows(); i++ {
		x := cvBoxMat.GetFloatAt(i, 0)
		y := cvBoxMat.GetFloatAt(i, 1)
		point := image.Pt(int(x), int(y))
		box = append(box, point)
	}

	cvBox := gocv.NewPointVectorFromPoints(box)
	cvDestination := gocv.NewPointVectorFromPoints(
		[]image.Point{
			{X: 0, Y: size.Y - 1},
			{X: 0, Y: 0},
			{X: size.X - 1, Y: 0},
			{X: size.X - 1, Y: size.Y - 1},
		},
	)
	cvTransformation := gocv.GetPerspectiveTransform(cvBox, cvDestination)

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
	cvPointsVector := gocv.NewPointsVectorFromPoints([][]image.Point{points})
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

func SaveRawImage(img gocv.Mat, basePath string) error {
	path := fmt.Sprintf("%s/raw.jpg", basePath)
	ok := gocv.IMWrite(path, img)
	if !ok {
		return fmt.Errorf("failed to write raw image to %s", path)
	}
	return nil
}

func SaveResults(basePath string, spots []ParkingSpot, predictions []float32) error {
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
	path := fmt.Sprintf("%s/result.json", basePath)
	err = os.WriteFile(path, resultJSON, 0644)
	if err != nil {
		return fmt.Errorf("failed to write result to %s: %w", path, err)
	}
	return nil
}

func SaveVisualizations(img gocv.Mat, basePath string, spots []ParkingSpot, predictions []float32) error {
	for i, prediction := range predictions {
		spot := spots[i]
		VisualizeSpotPrediction(&img, spot, prediction)
	}
	path := fmt.Sprintf("%s/visualization.jpg", basePath)
	ok := gocv.IMWrite(path, img)
	if !ok {
		return fmt.Errorf("failed to write visualization to %s", path)
	}
	return nil
}

func SavePredictions(img gocv.Mat, basePath string, spots []ParkingSpot, predictions []float32) error {
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", basePath, err)
	}
	if err := SaveRawImage(img, basePath); err != nil {
		return err
	}
	if err := SaveResults(basePath, spots, predictions); err != nil {
		return err
	}
	if err := SaveVisualizations(img, basePath, spots, predictions); err != nil {
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
