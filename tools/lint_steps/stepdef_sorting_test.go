package main_test

import (
	"testing"

	lintSteps "github.com/git-town/git-town/tools/lint_steps"
	"github.com/shoenig/test/must"
)

func TestStepDefSorting(t *testing.T) {
	t.Parallel()

	t.Run("FindStepDefinitions", func(t *testing.T) {
		t.Parallel()
		fileContent := "\n" +
			"func defineSteps(sc *godog.ScenarioContext) {\n" +
			"	sc.Step(`^a coworker clones the repository$`, func(ctx context.Context) {\n" +
			"		state := ctx.Value(keyScenarioState).(*ScenarioState)\n" +
			"	})\n" +
			"\n" +
			"	sc.Step(`^a folder \"([^\"]*)\"$`, func(ctx context.Context, name string) {\n" +
			"	})"
		have := lintSteps.FindStepDefinitions(fileContent)
		want := []lintSteps.StepDefinition{
			{
				Text: `^a coworker clones the repository$`,
				Line: 3,
			},
			{
				Text: `^a folder "([^"]*)"$`,
				Line: 7,
			},
		}
		must.Eq(t, want, have)
	})

	t.Run("FindUnsortedStepDefs", func(t *testing.T) {
		t.Parallel()
		stepDefs := []lintSteps.StepDefinition{
			{
				Text: `^a regex`,
				Line: 1,
			},
			{
				Text: `^c regex`,
				Line: 2,
			},
			{
				Text: `^b regex`,
				Line: 3,
			},
		}
		have := lintSteps.FindUnsortedStepDefs(stepDefs)
		want := []lintSteps.StepDefinition{
			{
				Text: `^b regex`,
				Line: 2,
			},
			{
				Text: `^c regex`,
				Line: 3,
			},
		}
		must.Eq(t, want, have)
	})

	t.Run("NormalizeForSort", func(t *testing.T) {
		t.Parallel()
		tests := map[string]string{
			"abc":                        "abc",
			`branch "([^"]+)" is active`: "branchisactive",
		}
		for give, want := range tests {
			have := lintSteps.NormalizeForSort(give)
			must.EqOp(t, want, have)
		}
	})
}
