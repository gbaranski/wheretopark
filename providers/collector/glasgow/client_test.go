package glasgow_test

import (
	"fmt"
	"testing"
	"wheretopark/providers/collector/glasgow"
)

func TestData(t *testing.T) {
	data, err := glasgow.GetData()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("data: %+v\n\n\n", data)
}
