package subshell_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/subshell"
	"github.com/stretchr/testify/assert"
)

func TestOutput(t *testing.T) {
	t.Parallel()
	t.Run(".ContainsText()", func(t *testing.T) {
		t.Parallel()
		output := subshell.Output{Raw: "one two three"}
		assert.True(t, output.ContainsText("two"), "should contain 'two'")
		assert.False(t, output.ContainsText("zonk"), "should not contain 'zonk'")
	})

	t.Run(".ContainsLine()", func(t *testing.T) {
		t.Parallel()
		output := subshell.Output{Raw: "one two"}
		assert.True(t, output.ContainsLine("one two"), `should contain "one two"`)
		assert.False(t, output.ContainsLine("hello"), `partial match should return false`)
		assert.False(t, output.ContainsLine("zonk"), `should not contain "zonk"`)
	})

	t.Run(".Lines()", func(t *testing.T) {
		t.Parallel()
		output := subshell.Output{Raw: "one\ntwo\nthree"}
		have := output.Lines()
		want := []string{"one", "two", "three"}
		assert.Equal(t, want, have)
	})
}
