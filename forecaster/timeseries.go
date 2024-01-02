package forecaster

import (
	wheretopark "wheretopark/go"
)

type Timeseries struct {
	ParkingLots map[wheretopark.ID]ParkingLot `json:"parkingLots"`
}
