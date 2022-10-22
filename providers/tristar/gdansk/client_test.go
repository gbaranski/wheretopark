package gdansk_test

import (
	"log"
	"testing"
	"wheretopark/providers/tristar/gdansk"
)

func TestMetadatas(t *testing.T) {
	metadatas, err := gdansk.Metadatas()
	if err != nil {
		t.Fatal(err)
	}
	log.Println(metadatas)
}

func TestStates(t *testing.T) {
	states, err := gdansk.States()
	if err != nil {
		t.Fatal(err)
	}
	log.Println(states)
}
