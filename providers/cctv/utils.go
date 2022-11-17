package cctv

import (
	"image"

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
