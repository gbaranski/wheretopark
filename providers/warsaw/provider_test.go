package warsaw_test

import (
	"testing"
	"wheretopark/go/tester"
	"wheretopark/providers/warsaw"
)

func TestWarsaw(t *testing.T) {
	tester.ExamineProvider(t, warsaw.New())
}
