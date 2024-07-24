package main_test

import (
	"testing"

	lintSteps "github.com/git-town/git-town/tools/lint_steps"
	"github.com/shoenig/test/must"
)

func TestFindStepDefinitions(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		give := "\n" +
			"func defineSteps(sc *godog.ScenarioContext) {\n" +
			"	sc.Step(`^a coworker clones the repository$`, func(ctx context.Context) {\n" +
			"		state := ctx.Value(keyScenarioState).(*ScenarioState)\n" +
			"	})\n" +
			"\n" +
			"	sc.Step(`^a folder \"([^\"]*)\"$`, func(ctx context.Context, name string) {\n" +
			"	})"
		have := lintSteps.FindStepDefinitions(give)
		want := []string{"xxx"}
		must.Eq(t, want, have)
	})
}
