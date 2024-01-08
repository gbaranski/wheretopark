package lacity_test

import (
	"testing"
	"wheretopark/go/tester"
	"wheretopark/providers/lacity"
)

func TestLACity(t *testing.T) {
	tester.ExamineProvider(t, lacity.New())
}
