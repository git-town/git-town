package runstate_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/stretchr/testify/assert"
)

func TestStepList(t *testing.T) {
	t.Run("Append", func(t *testing.T) {
		t.Run("append a single step", func(t *testing.T) {
			list := runstate.StepList{List: []steps.Step{&steps.AbortMergeStep{}}}
			list.Append(&steps.StashOpenChangesStep{})
			want := []steps.Step{&steps.AbortMergeStep{}, &steps.StashOpenChangesStep{}}
			assert.Equal(t, want, list.List)
		})
		t.Run("append multiple steps", func(t *testing.T) {
			list := runstate.StepList{List: []steps.Step{&steps.AbortMergeStep{}}}
			list.Append(&steps.AbortRebaseStep{}, &steps.StashOpenChangesStep{})
			want := []steps.Step{&steps.AbortMergeStep{}, &steps.AbortRebaseStep{}, &steps.StashOpenChangesStep{}}
			assert.Equal(t, want, list.List)
		})
		t.Run("append no steps", func(t *testing.T) {
			list := runstate.StepList{List: []steps.Step{}}
			list.Append()
			assert.Equal(t, []steps.Step{}, list.List)
		})
	})

	t.Run("Prepend", func(t *testing.T) {
		t.Run("prepend a single step", func(t *testing.T) {
			list := runstate.StepList{List: []steps.Step{&steps.AbortMergeStep{}}}
			list.Prepend(&steps.StashOpenChangesStep{})
			want := []steps.Step{&steps.StashOpenChangesStep{}, &steps.AbortMergeStep{}}
			assert.Equal(t, want, list.List)
		})
		t.Run("prepend multiple steps", func(t *testing.T) {
			list := runstate.StepList{List: []steps.Step{&steps.AbortMergeStep{}}}
			list.Prepend(&steps.AbortRebaseStep{}, &steps.StashOpenChangesStep{})
			want := []steps.Step{&steps.AbortRebaseStep{}, &steps.StashOpenChangesStep{}, &steps.AbortMergeStep{}}
			assert.Equal(t, want, list.List)
		})
		t.Run("prepend no steps", func(t *testing.T) {
			list := runstate.StepList{List: []steps.Step{}}
			list.Prepend()
			assert.Equal(t, []steps.Step{}, list.List)
		})
	})
}
