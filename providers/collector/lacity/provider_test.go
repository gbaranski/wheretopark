package lacity_test

import (
	"testing"
	"wheretopark/go/provider/sequential"
	"wheretopark/providers/collector/lacity"
)

func TestProvider(t *testing.T) {
	provider, err := lacity.NewProvider()
	if err != nil {
		t.Fatal(err)
	}
	sequential.ExamineProvider(t, provider)
}
