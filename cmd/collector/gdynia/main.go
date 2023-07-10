package main

import (
	wheretopark "wheretopark/go"
	"wheretopark/go/provider/sequential"
	"wheretopark/providers/collector/gdynia"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	provider, err := gdynia.NewProvider()
	if err != nil {
		panic(err)
	}
	lambda.Start(func() (map[wheretopark.ID]wheretopark.ParkingLot, error) {
		return sequential.GetParkingLots(provider)
	})
}
