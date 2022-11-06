package gdansk_test

import (
	"testing"
	"wheretopark/providers/tristar/gdansk"

	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	provider, err := gdansk.NewProvider(nil)
	if err != nil {
		t.Fatal(err)
	}
	metadata, err := provider.GetMetadata()
	if err != nil {
		t.Fatal(err)
	}

	state, err := provider.GetState()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(metadata), len(state))
}
