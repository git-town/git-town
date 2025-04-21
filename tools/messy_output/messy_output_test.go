package messyoutput_test

import (
	"os"
	"testing"

	messyOutput "github.com/git-town/git-town/tools/messyoutput"
	"github.com/shoenig/test/must"
)

func TestMessyOutput(t *testing.T) {
	t.Run("all okay", func(t *testing.T) {
		t.Parallel()
		text := `
Feature: test

  Scenario: one
	  First step
		Second step

	Scenario: two
	  First step
`
		must.NoError(t, os.WriteFile("test", []byte(text), 0o744))
		defer os.Remove("test")
		regex := messyOutput.CompileRegex()
		feature := messyOutput.ReadGherkinFile("test")
		have := messyOutput.FindScenarios(feature, "test", regex)
		want := []messyOutput.ScenarioInfo{}
		must.Eq(t, want, have)
	})
}
