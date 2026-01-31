package bytestream

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestNewlineDelineated(t *testing.T) {
	t.Parallel()

	t.Run("Sanitize", func(t *testing.T) {
		t.Parallel()
		t.Run("empty input", func(t *testing.T) {
			t.Parallel()
			give := NewlineDelineated([]byte{})
			have := give.Sanitize()
			want := Sanitized([]byte{})
			must.SliceEqOp(t, want, have)
		})
		t.Run("no secrets", func(t *testing.T) {
			t.Parallel()
			give := NewlineDelineated([]byte("hello world"))
			have := give.Sanitize()
			want := Sanitized([]byte("hello world"))
			must.SliceEqOp(t, want, have)
		})
		t.Run("deprecated Codeberg Token", func(t *testing.T) {
			t.Parallel()
			give := NewlineDelineated([]byte("hello\n\ngit-town.codeberg-token\n1234567890\n\nworld"))
			have := give.Sanitize()
			want := Sanitized([]byte("hello\n\ngit-town.codeberg-token\n(redacted)\n\nworld"))
			must.SliceEqOp(t, want, have)
		})
		t.Run("GitHub Token", func(t *testing.T) {
			t.Parallel()
			give := NewlineDelineated([]byte("hello\n\ngit-town.github-token\n1234567890\n\nworld"))
			have := give.Sanitize()
			want := Sanitized([]byte("hello\n\ngit-town.github-token\n(redacted)\n\nworld"))
			must.SliceEqOp(t, want, have)
		})
		t.Run("GitLab Token", func(t *testing.T) {
			t.Parallel()
			give := NewlineDelineated([]byte("hello\n\ngit-town.gitlab-token\n1234567890\n\nworld"))
			have := give.Sanitize()
			want := Sanitized([]byte("hello\n\ngit-town.gitlab-token\n(redacted)\n\nworld"))
			must.SliceEqOp(t, want, have)
		})
		t.Run("Forgejo Token", func(t *testing.T) {
			t.Parallel()
			give := NewlineDelineated([]byte("hello\n\ngit-town.forgejo-token\n1234567890\n\nworld"))
			have := give.Sanitize()
			want := Sanitized([]byte("hello\n\ngit-town.forgejo-token\n(redacted)\n\nworld"))
			must.SliceEqOp(t, want, have)
		})
		t.Run("Bitbucket Token", func(t *testing.T) {
			t.Parallel()
			give := NewlineDelineated([]byte("hello\n\ngit-town.bitbucket-app-password\n1234567890\n\nworld"))
			have := give.Sanitize()
			want := Sanitized([]byte("hello\n\ngit-town.bitbucket-app-password\n(redacted)\n\nworld"))
			must.SliceEqOp(t, want, have)
		})
		t.Run("Gitea Token", func(t *testing.T) {
			t.Parallel()
			give := NewlineDelineated([]byte("hello\n\ngit-town.gitea-token\n1234567890\n\nworld"))
			have := give.Sanitize()
			want := Sanitized([]byte("hello\n\ngit-town.gitea-token\n(redacted)\n\nworld"))
			must.SliceEqOp(t, want, have)
		})
		t.Run("user email", func(t *testing.T) {
			t.Parallel()
			give := NewlineDelineated([]byte("hello\n\nuser.email\nuser@example.com\n\nworld"))
			have := give.Sanitize()
			want := Sanitized([]byte("hello\n\nuser.email\n(redacted)\n\nworld"))
			must.SliceEqOp(t, want, have)
		})
	})
}
