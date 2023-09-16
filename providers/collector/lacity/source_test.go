package lacity_test

import (
	"testing"
	wheretopark "wheretopark/go"
	"wheretopark/providers/collector/lacity"
)

func TestLACity(t *testing.T) {
	wheretopark.ExamineSequentialSource(t, lacity.New())
}
