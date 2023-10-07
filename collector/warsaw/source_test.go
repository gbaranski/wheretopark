package warsaw_test

import (
	"testing"
	"wheretopark/collector/warsaw"
	wheretopark "wheretopark/go"
)

func TestWarsaw(t *testing.T) {
	wheretopark.ExamineSource(t, warsaw.New())
}
