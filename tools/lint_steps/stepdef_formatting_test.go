package main_test

import (
	"testing"

	lintSteps "github.com/git-town/git-town/tools/lint_steps"
	"github.com/shoenig/test/must"
)

func TestStepDefFormatting(t *testing.T) {
	t.Parallel()

	t.Run("CheckStepDefinitions", func(t *testing.T) {
		t.Parallel()
		give := "\n" +
			"func defineSteps(sc *godog.ScenarioContext) {\n" +
			"	sc.Step(`^a coworker clones the repository$`, func(ctx context.Context) {\n" +
			"		state := ctx.Value(keyScenarioState).(*ScenarioState)\n" +
			"	})\n" +
			"\n" +
			"	sc.Step(\"^a folder \"([^\"]*)\"$`, func(ctx context.Context, name string) {\n" +
			"	})" +
			"\n" +
			"	sc.Step('^a folder \"([^\"]*)\"$`, func(ctx context.Context, name string) {\n" +
			"	})"
		have := lintSteps.CheckStepDefinitions(give)
		want := []lintSteps.StepDefinition{
			{
				Line: 7,
				Text: `sc.Step("`,
			},
			{
				Line: 9,
				Text: "sc.Step('",
			},
		}
		must.Eq(t, want, have)
	})
}
