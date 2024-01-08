package gdansk_test

import (
	"testing"
	wheretopark "wheretopark/go"
	"wheretopark/providers/gdansk"
)

func TestGdansk(t *testing.T) {
	wheretopark.ExamineSource(t, gdansk.New())
}
