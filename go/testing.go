package wheretopark

import (
	"context"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func ExamineParkingLots(t *testing.T, parkingLots map[ID]ParkingLot) {
	assert.NotEmpty(t, parkingLots)
	for id, parkingLot := range parkingLots {
		t.Logf("%+v", parkingLot)
		if err := parkingLot.Validate(); err != nil {
			t.Fatalf("invalid parking lot %s : %s", id, err)
		}
	}
}

func ExamineSource(t *testing.T, src Source) {
	ch, err := src.ParkingLots(context.TODO())
	if err != nil {
		log.Fatal().Err(err).Msg("source fail")
	}

	parkingLots := make(map[ID]ParkingLot)
	for parkingLot := range ch {
		for id, lot := range parkingLot {
			parkingLots[id] = lot
		}
	}
	ExamineParkingLots(t, parkingLots)
}
