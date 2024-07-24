package main_test

import (
	"testing"

	lintSteps "github.com/git-town/git-town/tools/lint_steps"
	"github.com/shoenig/test/must"
)

func TestLintSteps(t *testing.T) {
	t.Parallel()
	t.Run("FindStepDefinitions", func(t *testing.T) {
		t.Parallel()
		give := "\n" +
			"func defineSteps(sc *godog.ScenarioContext) {\n" +
			"	sc.Step(`^a coworker clones the repository$`, func(ctx context.Context) {\n" +
			"		state := ctx.Value(keyScenarioState).(*ScenarioState)\n" +
			"	})\n" +
			"\n" +
			"	sc.Step(`^a folder \"([^\"]*)\"$`, func(ctx context.Context, name string) {\n" +
			"	})"
		have := lintSteps.FindStepDefinitions(give)
		want := []string{
			`^a coworker clones the repository$`,
			`^a folder "([^"]*)"$`,
		}
		must.Eq(t, want, have)
	})

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
		want := []string{
			`sc.Step("`,
			"sc.Step('",
		}
		must.Eq(t, want, have)
	})
}
