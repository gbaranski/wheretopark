package glasgow_test

import (
	"testing"
	"wheretopark/providers/collector/glasgow"
)

func TestProvider(t *testing.T) {
	provider, err := glasgow.NewProvider()
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
	for id := range state {
		_, exists := metadata[id]
		if !exists {
			t.Fatalf("missing metadata for %s", id)
		}
	}
}
