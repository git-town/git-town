package main_test

import (
	"testing"

	testsSorted "github.com/git-town/git-town/tools/tests_sorted"
	"github.com/shoenig/test/must"
)

const testPath = "test.go"

func TestTestsSorted(t *testing.T) {
	t.Parallel()

	t.Run("LintFile", func(t *testing.T) {
		t.Parallel()
		t.Run("unsorted subtests", func(t *testing.T) {
			t.Parallel()
			fileContents := `
package main

func TestF(t *testing.T) {
	t.Run("foo")
	t.Run("bar")
}
`
			want := `
test.go:4:1 unsorted subtests, expected order:

bar
foo

`[1:]

			got, err := testsSorted.LintFile(testPath, fileContents)

			must.NoError(t, err)
			must.EqOp(t, want, got.String())
		})
	})
}
