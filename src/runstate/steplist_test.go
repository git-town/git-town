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

	t.Run("AppendList", func(t *testing.T) {
		t.Run("append a populated list", func(t *testing.T) {
			list := runstate.StepList{List: []steps.Step{&steps.AbortMergeStep{}}}
			other := runstate.StepList{List: []steps.Step{&steps.StashOpenChangesStep{}}}
			list.AppendList(other)
			want := []steps.Step{&steps.AbortMergeStep{}, &steps.StashOpenChangesStep{}}
			assert.Equal(t, want, list.List)
		})
		t.Run("append an empty list", func(t *testing.T) {
			list := runstate.StepList{List: []steps.Step{&steps.AbortMergeStep{}}}
			other := runstate.StepList{List: []steps.Step{}}
			list.AppendList(other)
			assert.Equal(t, []steps.Step{&steps.AbortMergeStep{}}, list.List)
		})
	})

	t.Run("IsEmpty", func(t *testing.T) {
		t.Run("list is empty", func(t *testing.T) {
			list := runstate.StepList{List: []steps.Step{}}
			assert.True(t, list.IsEmpty())
		})
		t.Run("list is not empty", func(t *testing.T) {
			list := runstate.StepList{List: []steps.Step{&steps.AbortMergeStep{}}}
			assert.False(t, list.IsEmpty())
		})
	})

	t.Run("Peek", func(t *testing.T) {
		t.Run("populated list", func(t *testing.T) {
			list := runstate.StepList{List: []steps.Step{&steps.AbortMergeStep{}, &steps.StashOpenChangesStep{}}}
			have := list.Peek()
			assert.Equal(t, &steps.AbortMergeStep{}, have, "returns the first element of the list")
			wantList := runstate.StepList{List: []steps.Step{&steps.AbortMergeStep{}, &steps.StashOpenChangesStep{}}}
			assert.Equal(t, wantList, list, "does not modify the list")
		})
		t.Run("empty list", func(t *testing.T) {
			list := runstate.StepList{List: []steps.Step{}}
			have := list.Peek()
			assert.Equal(t, nil, have)
			wantList := runstate.StepList{List: []steps.Step{}}
			assert.Equal(t, wantList, list)
		})
	})

	t.Run("Pop", func(t *testing.T) {
		t.Run("populated list", func(t *testing.T) {
			list := runstate.StepList{List: []steps.Step{&steps.AbortMergeStep{}, &steps.StashOpenChangesStep{}}}
			have := list.Pop()
			assert.Equal(t, &steps.AbortMergeStep{}, have, "returns the first element of the list")
			wantList := runstate.StepList{List: []steps.Step{&steps.StashOpenChangesStep{}}}
			assert.Equal(t, wantList, list, "remotes the popped element from the list")
		})
		t.Run("empty list", func(t *testing.T) {
			list := runstate.StepList{List: []steps.Step{}}
			have := list.Pop()
			assert.Equal(t, nil, have, "returns nil")
			wantList := runstate.StepList{List: []steps.Step{}}
			assert.Equal(t, wantList, list, "remotes the popped element from the list")
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
