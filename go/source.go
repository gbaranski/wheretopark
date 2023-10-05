package wheretopark

import (
	"context"
)

type Source interface {
	ParkingLots(context.Context) (<-chan map[ID]ParkingLot, error)
}
