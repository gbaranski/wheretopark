package cctv

import (
	"fmt"
	"image"
	"image/color"

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
func (spot *ParkingSpot) CropOn(cvFrame gocv.Mat) gocv.Mat {
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

func (spot *ParkingSpot) VisualizeOn(img *gocv.Mat, prediction float32) {
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

func (parkingLot *ParkingLot) RunPredictions(model *Model, img gocv.Mat) []float32 {
	predictions := make([]float32, len(parkingLot.Spots))
	for i, spot := range parkingLot.Spots {
		croppedImage := spot.CropOn(img)
		defer croppedImage.Close()
		prediction := model.Predict(croppedImage)
		spot.VisualizeOn(&img, prediction)
		predictions[i] = prediction
	}
	return predictions
}
