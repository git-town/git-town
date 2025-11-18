package main_test

import (
	"testing"

	lintSteps "github.com/git-town/git-town/tools/lint_steps"
	"github.com/shoenig/test/must"
)

func TestCheckStepRegexAnchors(t *testing.T) {
	t.Parallel()

	t.Run("all steps properly anchored", func(t *testing.T) {
		t.Parallel()
		give := []lintSteps.StepDefinition{
			{Line: 1, Text: "^the first step$"},
			{Line: 2, Text: "^the second step$"},
			{Line: 3, Text: "^a step with (.*) capture$"},
		}
		have := lintSteps.CheckStepRegexAnchors(give)
		must.SliceEmpty(t, have)
	})

	t.Run("multiple unanchored steps", func(t *testing.T) {
		t.Parallel()
		give := []lintSteps.StepDefinition{
			{Line: 1, Text: "^properly anchored step$"},
			{Line: 2, Text: "missing start$"},
			{Line: 3, Text: "^missing end"},
			{Line: 4, Text: "missing both"},
		}
		have := lintSteps.CheckStepRegexAnchors(give)
		want := []lintSteps.StepDefinition{
			{Line: 2, Text: "missing start$"},
			{Line: 3, Text: "^missing end"},
			{Line: 4, Text: "missing both"},
		}
		must.Eq(t, want, have)
	})

	t.Run("step missing both anchors", func(t *testing.T) {
		t.Parallel()
		give := []lintSteps.StepDefinition{
			{Line: 1, Text: "^properly anchored step$"},
			{Line: 2, Text: "missing both anchors"},
		}
		have := lintSteps.CheckStepRegexAnchors(give)
		want := []lintSteps.StepDefinition{
			{Line: 2, Text: "missing both anchors"},
		}
		must.Eq(t, want, have)
	})

	t.Run("step missing end anchor", func(t *testing.T) {
		t.Parallel()
		give := []lintSteps.StepDefinition{
			{Line: 1, Text: "^properly anchored step$"},
			{Line: 2, Text: "^missing end anchor"},
		}
		have := lintSteps.CheckStepRegexAnchors(give)
		want := []lintSteps.StepDefinition{
			{Line: 2, Text: "^missing end anchor"},
		}
		must.Eq(t, want, have)
	})

	t.Run("step missing start anchor", func(t *testing.T) {
		t.Parallel()
		give := []lintSteps.StepDefinition{
			{Line: 1, Text: "^properly anchored step$"},
			{Line: 2, Text: "missing start anchor$"},
		}
		have := lintSteps.CheckStepRegexAnchors(give)
		want := []lintSteps.StepDefinition{
			{Line: 2, Text: "missing start anchor$"},
		}
		must.Eq(t, want, have)
	})
}
