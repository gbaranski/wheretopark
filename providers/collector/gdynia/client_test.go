package gdynia_test

import (
	"log"
	"testing"
	"wheretopark/providers/collector/gdansk"
)

func TestMetadatas(t *testing.T) {
	metadatas, err := gdansk.GetMetadata()
	if err != nil {
		t.Fatal(err)
	}
	log.Println(metadatas)
}

func TestStates(t *testing.T) {
	states, err := gdansk.GetState()
	if err != nil {
		t.Fatal(err)
	}
	log.Println(states)
}
