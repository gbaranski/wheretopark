package main

import (
	wheretopark "wheretopark/go"
	"wheretopark/providers/collector/glasgow"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	provider, err := glasgow.NewProvider()
	if err != nil {
		panic(err)
	}
	lambda.Start(func() (events.APIGatewayProxyResponse, error) {
		parkingLots, err := provider.GetParkingLots()
		return wheretopark.CreateGatewayProxyResponse(parkingLots, err), nil
	})
}
