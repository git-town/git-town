package main_test

import (
	"testing"

	formatUnittests "github.com/git-town/git-town/tools/format_unittests"
	"github.com/shoenig/test/must"
)

func TestFormatUnittests(t *testing.T) {
	t.Parallel()

	t.Run("FormatFileContent", func(t *testing.T) {
		t.Parallel()
		t.Run("top-level subtests", func(t *testing.T) {
			t.Parallel()
			give := `
package hosting_test

import (
	"code.gitea.io/sdk/gitea"
)

func TestNewGiteaConnector(t *testing.T) {
	t.Parallel()

	t.Run("top-level test 1", func(t *testing.T) {
		t.Parallel()
		give := 123
	})

	t.Run("top-level test 2", func(t *testing.T) {
		t.Parallel()
		give := 123
	})
}`
			want := `
package hosting_test

import (
	"code.gitea.io/sdk/gitea"
)

func TestNewGiteaConnector(t *testing.T) {
	t.Parallel()

	t.Run("top-level test 1", func(t *testing.T) {
		t.Parallel()
		give := 123
	})

	t.Run("top-level test 2", func(t *testing.T) {
		t.Parallel()
		give := 123
	})
}`
			have := formatUnittests.FormatFileContent(give)
			must.EqOp(t, want, have)
		})

		t.Run("nested subtests", func(t *testing.T) {
			t.Parallel()
			give := `
package hosting_test

import (
	"code.gitea.io/sdk/gitea"
)

func TestNewGiteaConnector(t *testing.T) {
	t.Parallel()

	t.Run("top-level test 1", func(t *testing.T) {
		t.Parallel()
		t.Run("nested test 1a", func(t *testing.T) {
			t.Parallel()
			give := 123
		})
		t.Run("nested test 1b", func(t *testing.T) {
			t.Parallel()
			give := 123
		})
	})

	t.Run("top-level test 2", func(t *testing.T) {
		t.Parallel()
		t.Run("nested test 2a", func(t *testing.T) {
			t.Parallel()
			give := 123
		})
		t.Run("nested test 2b", func(t *testing.T) {
			t.Parallel()
			give := 123
		})
	})
}`
			want := `
package hosting_test

import (
	"code.gitea.io/sdk/gitea"
)

func TestNewGiteaConnector(t *testing.T) {
	t.Parallel()

	t.Run("top-level test 1", func(t *testing.T) {
		t.Parallel()
		t.Run("nested test 1a", func(t *testing.T) {
			t.Parallel()
			give := 123
		})
		t.Run("nested test 1b", func(t *testing.T) {
			t.Parallel()
			give := 123
		})
	})

	t.Run("top-level test 2", func(t *testing.T) {
		t.Parallel()
		t.Run("nested test 2a", func(t *testing.T) {
			t.Parallel()
			give := 123
		})
		t.Run("nested test 2b", func(t *testing.T) {
			t.Parallel()
			give := 123
		})
	})
}`
			have := formatUnittests.FormatFileContent(give)
			must.EqOp(t, want, have)
		})

		t.Run("no subtests", func(t *testing.T) {
			t.Parallel()
			give := `
package hosting_test

import (
	"code.gitea.io/sdk/gitea"
)

func TestNewGiteaConnector(t *testing.T) {
	t.Parallel()
	give := "123"
}`
			want := `
package hosting_test

import (
	"code.gitea.io/sdk/gitea"
)

func TestNewGiteaConnector(t *testing.T) {
	t.Parallel()
	give := "123"
}`
			have := formatUnittests.FormatFileContent(give)
			must.EqOp(t, want, have)
		})
	})

	t.Run("IsGoTestFile", func(t *testing.T) {
		t.Parallel()
		tests := map[string]bool{
			"/one/two/three_test.go": true,
			"/one/two/three.go":      false,
			"/one/two_test/three.go": false,
		}
		for give, want := range tests {
			have := formatUnittests.IsGoTestFile(give)
			must.EqOp(t, want, have)
		}
	})

	t.Run("IsTopLevelRunLine", func(t *testing.T) {
		t.Parallel()
		tests := map[string]bool{
			"\tt.Run(\"HasLocalBranch\", func(t *testing.T) {":   true,
			"\t\tt.Run(\"HasLocalBranch\", func(t *testing.T) {": false,
		}
		for give, want := range tests {
			have := formatUnittests.IsTopLevelRunLine(give)
			must.EqOp(t, want, have)
		}
	})
}
