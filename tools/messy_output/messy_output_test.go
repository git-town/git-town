package messyoutput_test

import (
	"os"
	"testing"

	messyOutput "github.com/git-town/git-town/tools/messyoutput"
	"github.com/shoenig/test/must"
)

func TestMessyOutput(t *testing.T) {
	t.Run("ReadGherkinFile", func(t *testing.T) {
		t.Run("vanilla", func(t *testing.T) {
			text := `
Feature: test

  Scenario: one
	  First step
		Second step

	Scenario: two
	  First step
`[1:]
			must.NoError(t, os.WriteFile("test", []byte(text), 0o744))
			defer os.Remove("test")
			regex := messyOutput.CompileRegex()
			feature := messyOutput.ReadGherkinFile("test")
			have := messyOutput.FindScenarios(feature, "test", regex)
			want := []messyOutput.ScenarioInfo{
				{
					File:    "test",
					HasStep: false,
					HasTag:  false,
					Line:    3,
				},
				{
					File:    "test",
					HasStep: false,
					HasTag:  false,
					Line:    7,
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("feature has tag", func(t *testing.T) {
			text := `
@messyoutput
Feature: test

  Scenario: one
	  First step
		Second step

	Scenario: two
	  First step
`[1:]
			must.NoError(t, os.WriteFile("test", []byte(text), 0o744))
			defer os.Remove("test")
			regex := messyOutput.CompileRegex()
			feature := messyOutput.ReadGherkinFile("test")
			have := messyOutput.FindScenarios(feature, "test", regex)
			want := []messyOutput.ScenarioInfo{
				{
					File:    "test",
					HasStep: false,
					HasTag:  true,
					Line:    4,
				},
				{
					File:    "test",
					HasStep: false,
					HasTag:  true,
					Line:    8,
				},
			}
			must.Eq(t, want, have)
		})
	})
}
