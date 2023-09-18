package cctv

import (
	"bytes"
	"fmt"
	"image"
	"image/color"

	"github.com/rs/zerolog/log"
	ffmpeg "github.com/u2takey/ffmpeg-go"
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

func IsVacant(prediction float32) bool {
	if prediction > 0.5 {
		return true
	} else {
		return false
	}
}

func VisualizeSpotPrediction(img *gocv.Mat, spot ParkingSpot, prediction float32) {
	drawingColor := color.RGBA{
		A: 255,
	}
	if IsVacant(prediction) {
		drawingColor.G = 255
	} else {
		drawingColor.R = 255
	}
	points := spot.ImagePoints()
	cvPoints := gocv.NewPointVectorFromPoints(points)
	defer cvPoints.Close()
	cvPointsVector := gocv.NewPointsVectorFromPoints([][]image.Point{points})
	defer cvPointsVector.Close()
	cvMinAreaRect := gocv.MinAreaRect(cvPoints)
	gocv.Polylines(img, cvPointsVector, true, drawingColor, 2)
	gocv.PutText(
		img,
		fmt.Sprintf("%.2f", prediction),
		cvMinAreaRect.Center,
		gocv.FontHersheyDuplex,
		1,
		drawingColor,
		3,
	)
	gocv.PutText(
		img,
		spot.Type,
		cvMinAreaRect.Center.Add(image.Point{X: 0, Y: 16}),
		gocv.FontHersheyDuplex,
		0.5,
		color.RGBA{R: 0, G: 0, B: 255},
		2,
	)
}

func ExtractSpots(img gocv.Mat, spots []ParkingSpot) []gocv.Mat {
	images := make([]gocv.Mat, len(spots))
	for i, spot := range spots {
		images[i] = CropSpot(img, &spot)
	}
	return images
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func GetImageFromCamera(url string) (gocv.Mat, error) {
	buf := bytes.NewBuffer(nil)

	stream := ffmpeg.Input(url).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 1)}).
		Output("pipe:", ffmpeg.KwArgs{"frames:v": 1, "update": 1, "format": "image2"}).
		Silent(true)

	if log.Trace().Enabled() {
		stream = stream.WithOutput(buf, log.With().Str("level", "trace").Logger())
	} else {
		stream = stream.WithOutput(buf)
	}
	if err := stream.Run(); err != nil {
		return gocv.NewMat(), fmt.Errorf("unable to capture frame: %v", err)
	}
	img, err := gocv.IMDecode(buf.Bytes(), gocv.IMReadColor)
	if err != nil {
		return gocv.NewMat(), fmt.Errorf("unable to create image: %v", err)
	}
	if img.Empty() {
		return gocv.NewMat(), fmt.Errorf("empty image")
	}
	return img, nil
}
