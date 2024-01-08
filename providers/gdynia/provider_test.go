package gdynia_test

import (
	"testing"
	"wheretopark/go/tester"
	"wheretopark/providers/gdynia"
)

func TestGdynia(t *testing.T) {
	tester.ExamineProvider(t, gdynia.New())
}
