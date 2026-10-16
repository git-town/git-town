//nolint:paralleltest
package main_test

import (
	"os"
	"testing"

	messy "github.com/git-town/git-town/tools/messy_output"
	"github.com/shoenig/test/must"
)

func TestMessyOutput(t *testing.T) {

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
		scenarios := []messy.ScenarioInfo{
			{
				HasStep: false,
				HasTag:  true,
				Line:    4,
			},
			{
				HasStep: false,
				HasTag:  true,
				Line:    8,
			},
		}
		assertScenarios(t, text, scenarios)
		assertErrors(t, scenarios, []string{
			"test:4  unnecessary @messyoutput tag\n",
			"test:8  unnecessary @messyoutput tag\n",
		})
	})

	t.Run("no tags and steps", func(t *testing.T) {
		text := `
Feature: test

  Scenario: one
	  First step
		Second step

	Scenario: two
	  First step
`[1:]
		scenarios := []messy.ScenarioInfo{
			{
				HasStep: false,
				HasTag:  false,
				Line:    3,
			},
			{
				HasStep: false,
				HasTag:  false,
				Line:    7,
			},
		}
		assertScenarios(t, text, scenarios)
		assertErrors(t, scenarios, []string{})
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
		scenarios := []messy.ScenarioInfo{
			{
				HasStep: false,
				HasTag:  true,
				Line:    4,
			},
			{
				HasStep: false,
				HasTag:  false,
				Line:    8,
			},
		}
		assertScenarios(t, text, scenarios)
		assertErrors(t, scenarios, []string{
			"test:4  unnecessary @messyoutput tag\n",
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
		scenarios := []messy.ScenarioInfo{
			{
				HasStep: true,
				HasTag:  true,
				Line:    4,
			},
			{
				HasStep: false,
				HasTag:  false,
				Line:    8,
			},
		}
		assertScenarios(t, text, scenarios)
		assertErrors(t, scenarios, []string{})
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
		scenarios := []messy.ScenarioInfo{
			{
				HasStep: true,
				HasTag:  false,
				Line:    3,
			},
			{
				HasStep: false,
				HasTag:  false,
				Line:    7,
			},
		}
		assertScenarios(t, text, scenarios)
		assertErrors(t, scenarios, []string{
			"test:3  missing @messyoutput tag\n",
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
		scenarios := []messy.ScenarioInfo{
			{
				HasStep: true,
				HasTag:  false,
				Line:    3,
			},
			{
				HasStep: false,
				HasTag:  false,
				Line:    7,
			},
		}
		assertScenarios(t, text, scenarios)
		assertErrors(t, scenarios, []string{
			"test:3  missing @messyoutput tag\n",
		})
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
		scenarios := []messy.ScenarioInfo{
			{
				HasStep: true,
				HasTag:  true,
				Line:    4,
			},
			{
				HasStep: true,
				HasTag:  true,
				Line:    7,
			},
		}
		assertScenarios(t, text, scenarios)
		assertErrors(t, scenarios, []string{})
	})

	t.Run("the feature has the tag, only one scenarios has the step", func(t *testing.T) {
		text := `
@messyoutput
Feature: test

  Scenario: one
		And I run "foo" and enter into the dialogs:

	Scenario: two
		And another step here
`[1:]
		scenarios := []messy.ScenarioInfo{
			{
				HasStep: true,
				HasTag:  true,
				Line:    4,
			},
			{
				HasStep: false,
				HasTag:  true,
				Line:    7,
			},
		}
		assertScenarios(t, text, scenarios)
		assertErrors(t, scenarios, []string{
			"test:7  unnecessary @messyoutput tag\n",
		})
	})
}

func assertErrors(t *testing.T, scenarios []messy.ScenarioInfo, wantErrors []string) {
	t.Helper()
	haveErrors := messy.AnalyzeScenarios("test", scenarios)
	must.Eq(t, wantErrors, haveErrors)
}

func assertScenarios(t *testing.T, text string, wantScenarios []messy.ScenarioInfo) {
	t.Helper()
	must.NoError(t, os.WriteFile("test", []byte(text), 0o600))
	defer os.Remove("test")
	regex := messy.CompileRegex()
	feature := messy.ReadGherkinFile("test")
	haveScenarios := messy.FindScenarios(feature, regex)
	must.Eq(t, wantScenarios, haveScenarios)
}
