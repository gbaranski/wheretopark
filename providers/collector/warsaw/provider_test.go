package warsaw_test

import (
	"testing"
	"wheretopark/go/provider/simple"
	"wheretopark/providers/collector/warsaw"
)

func TestProvider(t *testing.T) {
	provider, err := warsaw.NewProvider()
	if err != nil {
		t.Fatal(err)
	}
	simple.ExamineProvider(t, provider)
}
