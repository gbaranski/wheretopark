package glasgow_test

import (
	"testing"
	"wheretopark/go/tester"
	"wheretopark/providers/glasgow"
)

func TestGlasgow(t *testing.T) {
	tester.ExamineProvider(t, glasgow.New())
}
