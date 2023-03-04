package gdynia_test

import (
	"log"
	"testing"
	"wheretopark/providers/collector/gdynia"
)

func TestMetadatas(t *testing.T) {
	metadatas, err := gdynia.GetMetadata()
	if err != nil {
		t.Fatal(err)
	}
	log.Println(metadatas)
}

func TestStates(t *testing.T) {
	states, err := gdynia.GetState()
	if err != nil {
		t.Fatal(err)
	}
	log.Println(states)
}
