package glasgow_test

import (
	"testing"
	wheretopark "wheretopark/go"
	"wheretopark/providers/glasgow"
)

func TestGlasgow(t *testing.T) {
	wheretopark.ExamineSource(t, glasgow.New())
}
