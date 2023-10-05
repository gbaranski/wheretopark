package lacity_test

import (
	"testing"
	"wheretopark/collector/lacity"
	wheretopark "wheretopark/go"
)

func TestLACity(t *testing.T) {
	wheretopark.ExamineSource(t, lacity.New())
}
