package main_test

import (
	"testing"

	lintSteps "github.com/git-town/git-town/tools/lint_steps"
	"github.com/shoenig/test/must"
)

func TestStepDefUsage(t *testing.T) {
	t.Parallel()

	t.Run("FindUnusedStepDefs", func(t *testing.T) {
		t.Parallel()

		t.Run("no unused steps", func(t *testing.T) {
			t.Parallel()
			have := lintSteps.FindAllUnusedStepDefs()
			want := []string{}
			must.Eq(t, want, have)
		})

		t.Run("some unused steps", func(t *testing.T) {
			t.Parallel()
		})
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
