package warsaw_test

import (
	"testing"
	"wheretopark/providers/collector/warsaw"

	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	provider, err := warsaw.NewProvider()
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
