package glasgow_test

import (
	"testing"
	"wheretopark/collector/glasgow"
	wheretopark "wheretopark/go"
)

func TestGlasgow(t *testing.T) {
	wheretopark.ExamineSource(t, glasgow.New())
}
