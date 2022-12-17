package gdynia_test

import (
	"testing"
	"wheretopark/providers/collector/gdynia"

	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	provider, err := gdynia.NewProvider(nil)
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
