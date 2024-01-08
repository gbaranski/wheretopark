package lacity_test

import (
	"testing"
	wheretopark "wheretopark/go"
	"wheretopark/providers/lacity"
)

func TestLACity(t *testing.T) {
	wheretopark.ExamineSource(t, lacity.New())
}
