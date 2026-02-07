package pagination

import (
	"testing"
)

func TestSpan_LimitSize(t *testing.T) {
	t.Run("less", func(t *testing.T) {
		span := NewFirstSpan(10)
		if span.LimitSize(3) != 3 {
			t.Fatal("should decrease size")
		}
	})

	t.Run("more", func(t *testing.T) {
		span := NewFirstSpan(10)
		if span.LimitSize(13) != 10 {
			t.Fatal("should not decrease size")
		}
	})
}
