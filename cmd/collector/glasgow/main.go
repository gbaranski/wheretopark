package main

import (
	"wheretopark/providers/collector/glasgow"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	provider, err := glasgow.NewProvider()
	if err != nil {
		panic(err)
	}
	lambda.Start(provider.GetParkingLots)
}
