package libwebring

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestRingSurroundingLinks(t *testing.T) {
	t.Run("0", func(t *testing.T) {
		ring := Ring{}

		l, r := ring.SurroundingIndex(0)
		assert.Equal(t, Link{}, l)
		assert.Equal(t, Link{}, r)
	})

	t.Run("1", func(t *testing.T) {
		ring := Ring{
			{"a", "https://a"},
		}

		l, r := ring.SurroundingIndex(0)
		assert.Equal(t, Link{"a", "https://a"}, l)
		assert.Equal(t, Link{"a", "https://a"}, r)
	})

	t.Run("2", func(t *testing.T) {
		ring := Ring{
			{"a", "https://a"},
			{"b", "https://b"},
		}

		l, r := ring.SurroundingIndex(0)
		assert.Equal(t, Link{"b", "https://b"}, l)
		assert.Equal(t, Link{"b", "https://b"}, r)
	})

	t.Run("3", func(t *testing.T) {
		ring := Ring{
			{"a", "https://a"},
			{"b", "https://b"},
			{"c", "https://c"},
		}

		l, r := ring.SurroundingIndex(0)
		assert.Equal(t, Link{"c", "https://c"}, l)
		assert.Equal(t, Link{"b", "https://b"}, r)
	})
}
