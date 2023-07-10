package main

import (
	wheretopark "wheretopark/go"
	"wheretopark/go/provider/sequential"
	"wheretopark/providers/collector/gdansk"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	provider, err := gdansk.NewProvider()
	if err != nil {
		panic(err)
	}
	lambda.Start(func() (map[wheretopark.ID]wheretopark.ParkingLot, error) {
		return sequential.GetParkingLots(provider)
	})
}
