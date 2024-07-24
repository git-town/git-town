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
				Line: 2,
			},
			{
				Text: `^a folder "([^"]*)"$`,
				Line: 6,
			},
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

	t.Run("FindUsedStepsIn", func(t *testing.T) {
		t.Parallel()
		fileContent := "\n" +
			"Feature: test\n" +
			"\n" +
			"  Background:\n" +
			"		Given step one\n" +
			"	  And step two\n" +
			"	  When step three\n" +
			"\n" +
			"  Scenario: result\n" +
			"	  Then step four\n" +
			"	  And step five\n" +
			"\n" +
			"  Scenario: undo\n" +
			"	  When step six\n" +
			"	  Then step seven\n" +
			"	  And step eight\n"
		have := lintSteps.FindUsedStepsIn(fileContent)
		want := []string{
			"step one",
			"step two",
			"step three",
			"step four",
			"step five",
			"step six",
			"step seven",
			"step eight",
		}
		must.Eq(t, want, have)
	})
}
