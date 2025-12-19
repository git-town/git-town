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
				Line: 3,
				Text: `^a coworker clones the repository$`,
			},
			{
				Line: 7,
				Text: `^a folder "([^"]*)"$`,
			},
		}
		must.Eq(t, want, have)
	})
}
