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

func ExamineProvider(t *testing.T, provider Provider) {
	sources := provider.Sources()
	for identifier, source := range sources {
		ctx := log.With().Str("source", identifier).Logger().WithContext(context.TODO())
		parkingLots, err := source.ParkingLots(ctx)
		if err != nil {
			t.Fatalf("source %s fail: %s", identifier, err)
		}
		ExamineParkingLots(t, parkingLots)
	}
}

func ExamineSource(t *testing.T, src Source) {
	parkingLots, err := src.ParkingLots(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	ExamineParkingLots(t, parkingLots)
}

func ExamineSequentialSource(t *testing.T, src SequentialSource) {
	proxy := NewSequentialSourceProxy(src)
	ExamineSource(t, proxy)
}
