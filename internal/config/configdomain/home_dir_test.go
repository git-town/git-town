package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestHomeDir_Join(t *testing.T) {
	t.Parallel()

	t.Run("no elements", func(t *testing.T) {
		t.Parallel()
		homeDir := configdomain.ConfigDirRepo("/home/user")
		have := homeDir.Join()
		want := "/home/user"
		must.Eq(t, want, have)
	})

	t.Run("single element", func(t *testing.T) {
		t.Parallel()
		homeDir := configdomain.ConfigDirRepo("/home/user")
		have := homeDir.Join("documents")
		want := "/home/user/documents"
		must.Eq(t, want, have)
	})

	t.Run("multiple elements", func(t *testing.T) {
		t.Parallel()
		homeDir := configdomain.ConfigDirRepo("/home/user")
		have := homeDir.Join("documents", "projects", "git-town")
		want := "/home/user/documents/projects/git-town"
		must.Eq(t, want, have)
	})

	t.Run("elements with path separators", func(t *testing.T) {
		t.Parallel()
		homeDir := configdomain.ConfigDirRepo("/home/user")
		have := homeDir.Join("documents/projects", "git-town")
		want := "/home/user/documents/projects/git-town"
		must.Eq(t, want, have)
	})

	t.Run("empty string elements", func(t *testing.T) {
		t.Parallel()
		homeDir := configdomain.ConfigDirRepo("/home/user")
		have := homeDir.Join("", "documents", "")
		want := "/home/user/documents"
		must.Eq(t, want, have)
	})

	t.Run("absolute path element", func(t *testing.T) {
		t.Parallel()
		homeDir := configdomain.ConfigDirRepo("/home/user")
		have := homeDir.Join("/tmp", "file")
		want := "/home/user/tmp/file"
		must.Eq(t, want, have)
	})
}
