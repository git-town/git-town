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
			scenarios := assert_scenarios(t, text, []messyOutput.ScenarioInfo{
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
			})
			assert_errors(t, scenarios, []string{})
		})

		t.Run("feature has tag but no steps", func(t *testing.T) {
			text := `
@messyoutput
Feature: test

  Scenario: one
	  First step
		Second step

	Scenario: two
	  First step
`[1:]
			scenarios := assert_scenarios(t, text, []messyOutput.ScenarioInfo{
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
			})
			assert_errors(t, scenarios, []string{
				"test:4  unnecessary tag\n",
				"test:8  unnecessary tag\n",
			})
		})

		t.Run("one scenario has a tag but no steps", func(t *testing.T) {
			text := `
Feature: test

	@messyoutput
  Scenario: one
	  First step
		Second step

	Scenario: two
	  First step
`[1:]
			scenarios := assert_scenarios(t, text, []messyOutput.ScenarioInfo{
				{
					File:    "test",
					HasStep: false,
					HasTag:  true,
					Line:    4,
				},
				{
					File:    "test",
					HasStep: false,
					HasTag:  false,
					Line:    8,
				},
			})
			assert_errors(t, scenarios, []string{
				"test:4  unnecessary tag\n",
			})
		})

		t.Run("one scenario has the step in singular but no tag", func(t *testing.T) {
			text := `
Feature: test

  Scenario: one
	  First step
		And I run "foo" and enter into the dialog:

	Scenario: two
	  First step
`[1:]
			scenarios := assert_scenarios(t, text, []messyOutput.ScenarioInfo{
				{
					File:    "test",
					HasStep: true,
					HasTag:  false,
					Line:    3,
				},
				{
					File:    "test",
					HasStep: false,
					HasTag:  false,
					Line:    7,
				},
			})
			assert_errors(t, scenarios, []string{
				"test:3  missing tag\n",
			})
		})

		t.Run("one scenario has the step in plural", func(t *testing.T) {
			text := `
Feature: test

  Scenario: one
	  First step
		And I run "foo" and enter into the dialogs:

	Scenario: two
	  First step
`[1:]
			scenarios := assert_scenarios(t, text, []messyOutput.ScenarioInfo{
				{
					File:    "test",
					HasStep: true,
					HasTag:  false,
					Line:    3,
				},
				{
					File:    "test",
					HasStep: false,
					HasTag:  false,
					Line:    7,
				},
			})
			assert_errors(t, scenarios, []string{
				"test:3  missing tag\n",
			})
		})

		t.Run("one scenario has the step and a tag", func(t *testing.T) {
			text := `
Feature: test

  @messyoutput
  Scenario: one
	  First step
		And I run "foo" and enter into the dialogs:

	Scenario: two
	  First step
`[1:]
			scenarios := assert_scenarios(t, text, []messyOutput.ScenarioInfo{
				{
					File:    "test",
					HasStep: true,
					HasTag:  true,
					Line:    4,
				},
				{
					File:    "test",
					HasStep: false,
					HasTag:  false,
					Line:    8,
				},
			})
			assert_errors(t, scenarios, []string{})
		})

		t.Run("the feature has the tag, both scenarios have the step", func(t *testing.T) {
			text := `
@messyoutput
Feature: test

  Scenario: one
		And I run "foo" and enter into the dialogs:

	Scenario: two
		And I run "foo" and enter into the dialogs:
`[1:]
			scenarios := assert_scenarios(t, text, []messyOutput.ScenarioInfo{
				{
					File:    "test",
					HasStep: true,
					HasTag:  true,
					Line:    4,
				},
				{
					File:    "test",
					HasStep: true,
					HasTag:  true,
					Line:    7,
				},
			})
			assert_errors(t, scenarios, []string{})
		})

		t.Run("the feature has the tag, both scenarios have the step", func(t *testing.T) {
			text := `
@messyoutput
Feature: test

  Scenario: one
		And I run "foo" and enter into the dialogs:

	Scenario: two
		And another step here
`[1:]
			scenarios := assert_scenarios(t, text, []messyOutput.ScenarioInfo{
				{
					File:    "test",
					HasStep: true,
					HasTag:  true,
					Line:    4,
				},
				{
					File:    "test",
					HasStep: false,
					HasTag:  true,
					Line:    7,
				},
			})
			assert_errors(t, scenarios, []string{
				"test:7  unnecessary tag\n",
			})
		})
	})
}

func assert_scenarios(t *testing.T, text string, wantScenarios []messyOutput.ScenarioInfo) []messyOutput.ScenarioInfo {
	must.NoError(t, os.WriteFile("test", []byte(text), 0o744))
	defer os.Remove("test")
	regex := messyOutput.CompileRegex()
	feature := messyOutput.ReadGherkinFile("test")
	haveScenarios := messyOutput.FindScenarios(feature, "test", regex)
	must.Eq(t, wantScenarios, haveScenarios)
	return haveScenarios
}

func assert_errors(t *testing.T, scenarios []messyOutput.ScenarioInfo, wantErrors []string) {
	t.Helper()
	haveErrors := messyOutput.AnalyzeScenarios(scenarios)
	must.Eq(t, wantErrors, haveErrors)
}
