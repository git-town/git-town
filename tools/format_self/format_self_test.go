package main_test

import (
	"testing"

	formatSelf "github.com/git-town/git-town/tools/format_self"
	"github.com/shoenig/test/must"
)

func TestFormatSelf(t *testing.T) {
	t.Parallel()

	t.Run("FormatFileContent", func(t *testing.T) {
		t.Parallel()
		t.Run("unformatted, non-pointer receiver", func(t *testing.T) {
			t.Parallel()
			give := `
package main
type Foo struct{}
func (f Foo) Bar() {
	// ...
}`
			want := `
package main
type Foo struct{}
func (self Foo) Bar() {
	// ...
}`
			have := formatSelf.FormatFileContent(give)
			must.EqOp(t, want, have)
		})
		t.Run("unformatted, pointer receiver", func(t *testing.T) {
			t.Parallel()
			give := `
package main
type Foo struct{}
func (f *Foo) Bar() {
	fmt.Println("")
}`
			have := formatSelf.FormatFileContent(give)
			want := `
package main
type Foo struct{}
func (self *Foo) Bar() {
	fmt.Println("")
}`
			must.EqOp(t, want, have)
		})
		t.Run("unformatted, generic method", func(t *testing.T) {
			t.Parallel()
			give := `
package main
type WithPrevious struct{}
func (c *WithPrevious[T]) Initialized() bool {
        return c.initialized
}`
			have := formatSelf.FormatFileContent(give)
			want := `
package main
type WithPrevious struct{}
func (self *WithPrevious[T]) Initialized() bool {
        return c.initialized
}`
			must.EqOp(t, want, have)
		})
		t.Run("already formatted", func(t *testing.T) {
			t.Parallel()
			give := `
			package main
			type Foo struct{}
			func (self Foo) Bar() {
				fmt.Println("")
			}
			`
			have := formatSelf.FormatFileContent(give)
			must.EqOp(t, give, have)
		})
	})

	t.Run("FormatLine", func(t *testing.T) {
		t.Parallel()
		tests := map[string]string{
			"func (bcs *BackendCommands) CommentOutSquashCommitMessage(prefix string) error {": "func (self *BackendCommands) CommentOutSquashCommitMessage(prefix string) error {",
			"func (c *Counter) Count() int {":                                                  "func (self *Counter) Count() int {",
			"	if err != nil {":                                                                 "	if err != nil {",
		}
		for give, want := range tests {
			have := formatSelf.FormatLine(give)
			must.EqOp(t, want, have)
		}
	})

	t.Run("IsGoFile", func(t *testing.T) {
		t.Parallel()
		tests := map[string]bool{
			"/foo/bar.go":      true,
			"/foo/bar_test.go": false,
			"/foo/bar.md":      false,
			"/foo/bar.go.md":   false,
		}
		for give, want := range tests {
			have := formatSelf.IsGoFile(give)
			must.EqOp(t, want, have)
		}
	})
}
