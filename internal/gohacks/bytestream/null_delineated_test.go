package bytestream_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/gohacks/bytestream"
	"github.com/shoenig/test/must"
)

func TestNullDelineated(t *testing.T) {
	t.Parallel()

	t.Run("ToNewlines", func(t *testing.T) {
		t.Parallel()

		t.Run("empty input", func(t *testing.T) {
			t.Parallel()
			give := bytestream.NullDelineated([]byte{})
			have := give.ToNewlines()
			want := bytestream.NewlineDelineated([]byte{})
			must.SliceEqOp(t, want, have)
		})

		t.Run("no null bytes", func(t *testing.T) {
			t.Parallel()
			give := bytestream.NullDelineated([]byte("hello world"))
			have := give.ToNewlines()
			want := bytestream.NewlineDelineated([]byte("hello world"))
			must.SliceEqOp(t, want, have)
		})

		t.Run("single null byte", func(t *testing.T) {
			t.Parallel()
			give := bytestream.NullDelineated([]byte{'h', 'e', 'l', 'l', 'o', 0x00, 'w', 'o', 'r', 'l', 'd'})
			have := give.ToNewlines()
			want := bytestream.NewlineDelineated([]byte{'h', 'e', 'l', 'l', 'o', '\n', '\n', 'w', 'o', 'r', 'l', 'd'})
			must.SliceEqOp(t, want, have)
		})

		t.Run("multiple null bytes", func(t *testing.T) {
			t.Parallel()
			give := bytestream.NullDelineated([]byte{'a', 0x00, 'b', 0x00, 'c'})
			have := give.ToNewlines()
			want := bytestream.NewlineDelineated([]byte{'a', '\n', '\n', 'b', '\n', '\n', 'c'})
			must.SliceEqOp(t, want, have)
		})

		t.Run("null byte at beginning", func(t *testing.T) {
			t.Parallel()
			give := bytestream.NullDelineated([]byte{0x00, 'h', 'e', 'l', 'l', 'o'})
			have := give.ToNewlines()
			want := bytestream.NewlineDelineated([]byte{'\n', '\n', 'h', 'e', 'l', 'l', 'o'})
			must.SliceEqOp(t, want, have)
		})

		t.Run("null byte at end", func(t *testing.T) {
			t.Parallel()
			give := bytestream.NullDelineated([]byte{'h', 'e', 'l', 'l', 'o', 0x00})
			have := give.ToNewlines()
			want := bytestream.NewlineDelineated([]byte{'h', 'e', 'l', 'l', 'o', '\n', '\n'})
			must.SliceEqOp(t, want, have)
		})

		t.Run("only null bytes", func(t *testing.T) {
			t.Parallel()
			give := bytestream.NullDelineated([]byte{0x00, 0x00, 0x00})
			have := give.ToNewlines()
			want := bytestream.NewlineDelineated([]byte{'\n', '\n', '\n', '\n', '\n', '\n'})
			must.SliceEqOp(t, want, have)
		})

		t.Run("consecutive null bytes", func(t *testing.T) {
			t.Parallel()
			give := bytestream.NullDelineated([]byte{'a', 0x00, 0x00, 'b'})
			have := give.ToNewlines()
			want := bytestream.NewlineDelineated([]byte{'a', '\n', '\n', '\n', '\n', 'b'})
			must.SliceEqOp(t, want, have)
		})
	})
}
