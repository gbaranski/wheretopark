package main

import (
	"wheretopark/providers/collector/poznan"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	provider, err := poznan.NewProvider()
	if err != nil {
		panic(err)
	}
	lambda.Start(provider.GetParkingLots)
}
