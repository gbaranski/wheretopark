package simple

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExamineProvider(t *testing.T, provider Provider) {
	parkingLots, err := provider.GetParkingLots()
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, parkingLots)

	for id, parkingLot := range parkingLots {
		t.Logf("%+v", parkingLot)
		if err := parkingLot.Validate(); err != nil {
			t.Fatalf("invalid parking lot %s : %e", id, err)
		}
	}
}
