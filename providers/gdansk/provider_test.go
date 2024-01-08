package gdansk_test

import (
	"testing"
	"wheretopark/go/tester"
	"wheretopark/providers/gdansk"
)

func TestGdansk(t *testing.T) {
	tester.ExamineProvider(t, gdansk.New())
}
