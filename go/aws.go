package wheretopark

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func createGatewayProxyErrorResponse(err error, status int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode:        status,
		Headers:           map[string]string{},
		MultiValueHeaders: map[string][]string{},
		Body:              err.Error(),
		IsBase64Encoded:   false,
	}
}

func CreateGatewayProxyResponse(parkingLots map[ID]ParkingLot, providerError error) events.APIGatewayProxyResponse {
	if providerError != nil {
		return createGatewayProxyErrorResponse(providerError, http.StatusBadGateway)
	}
	json, err := json.Marshal(parkingLots)
	if err != nil {
		return createGatewayProxyErrorResponse(err, http.StatusInternalServerError)
	}
	return events.APIGatewayProxyResponse{
		StatusCode:        200,
		Headers:           map[string]string{},
		MultiValueHeaders: map[string][]string{},
		Body:              string(json),
		IsBase64Encoded:   false,
	}

}
