package gdansk_test

import (
	"testing"
	"wheretopark/go/provider/sequential"
	"wheretopark/providers/collector/gdansk"
)

func TestProvider(t *testing.T) {
	provider, err := gdansk.NewProvider()
	if err != nil {
		t.Fatal(err)
	}
	sequential.ExamineProvider(t, provider)
}
