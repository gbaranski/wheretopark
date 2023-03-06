package poznan_test

import (
	"log"
	"testing"
	"wheretopark/providers/collector/poznan"
)

func TestData(t *testing.T) {
	data, err := poznan.GetData()
	if err != nil {
		t.Fatal(err)
	}
	log.Println(data)
}
