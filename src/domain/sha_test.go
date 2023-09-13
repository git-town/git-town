package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/test/asserts"
	"github.com/stretchr/testify/assert"
)

func TestSHA(t *testing.T) {
	t.Parallel()
	t.Run("IsEmpty", func(t *testing.T) {
		t.Parallel()
		t.Run("is empty", func(t *testing.T) {
			t.Parallel()
			sha := domain.SHA{}
			assert.True(t, sha.IsEmpty())
		})
		t.Run("is not empty", func(t *testing.T) {
			t.Parallel()
			sha := domain.NewSHA("123456")
			assert.False(t, sha.IsEmpty())
		})
	})

	t.Run("NewSHA and String", func(t *testing.T) {
		t.Parallel()
		t.Run("allows lowercase hex characters", func(t *testing.T) {
			t.Parallel()
			text := "1234567890abcdef"
			sha := domain.NewSHA(text)
			assert.Equal(t, text, sha.String())
		})
		t.Run("does not allow empty values", func(t *testing.T) {
			t.Parallel()
			defer asserts.Paniced(t)
			domain.NewSHA("")
		})
		t.Run("does not allow spaces", func(t *testing.T) {
			t.Parallel()
			defer asserts.Paniced(t)
			domain.NewSHA("abc def")
		})
		t.Run("does not allow uppercase characters", func(t *testing.T) {
			t.Parallel()
			defer asserts.Paniced(t)
			domain.NewSHA("ABCDEF")
		})
		t.Run("does not allow non-hex characters", func(t *testing.T) {
			t.Parallel()
			defer asserts.Paniced(t)
			domain.NewSHA("abcdefg")
		})
	})

	t.Run("TruncateTo", func(t *testing.T) {
		t.Parallel()
		t.Run("SHA is longer than the new length", func(t *testing.T) {
			t.Parallel()
			sha := domain.NewSHA("123456789abcdef")
			have := sha.TruncateTo(8)
			want := domain.NewSHA("12345678")
			assert.Equal(t, want, have)
		})
		t.Run("SHA is shorter than the new length", func(t *testing.T) {
			t.Parallel()
			sha := domain.NewSHA("123456789")
			have := sha.TruncateTo(12)
			want := domain.NewSHA("123456789")
			assert.Equal(t, want, have)
		})
	})
}
