package sequential

import (
	"testing"
	wheretopark "wheretopark/go"

	"github.com/stretchr/testify/assert"
)

func ExamineProvider(t *testing.T, provider Provider) {
	metadata, err := provider.GetMetadatas()
	if err != nil {
		t.Fatal(err)
	}

	state, err := provider.GetStates()
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, metadata)
	assert.NotEmpty(t, state)

	for id := range metadata {
		_, exists := state[id]
		if !exists {
			delete(metadata, id)
		}
	}
	for id := range state {
		_, exists := metadata[id]
		if !exists {
			delete(state, id)
		}
	}

	assert.Equal(t, len(metadata), len(state))

	for id, metadata := range metadata {
		state := state[id]
		parkingLot := wheretopark.ParkingLot{
			State:    state,
			Metadata: metadata,
		}
		t.Logf("%+v", parkingLot)
		if err := parkingLot.Validate(); err != nil {
			t.Fatalf("invalid parking lot %s : %e", id, err)
		}
	}
}
