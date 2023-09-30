package runstate_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/stretchr/testify/assert"
)

func TestStepList(t *testing.T) {
	t.Parallel()

	t.Run("Append", func(t *testing.T) {
		t.Parallel()
		t.Run("append a single step", func(t *testing.T) {
			t.Parallel()
			list := runstate.StepList{List: []steps.Step{&steps.AbortMergeStep{}}}
			list.Append(&steps.StashOpenChangesStep{})
			want := []steps.Step{&steps.AbortMergeStep{}, &steps.StashOpenChangesStep{}}
			assert.Equal(t, want, list.List)
		})
		t.Run("append multiple steps", func(t *testing.T) {
			t.Parallel()
			list := runstate.StepList{List: []steps.Step{&steps.AbortMergeStep{}}}
			list.Append(&steps.AbortRebaseStep{}, &steps.StashOpenChangesStep{})
			want := []steps.Step{&steps.AbortMergeStep{}, &steps.AbortRebaseStep{}, &steps.StashOpenChangesStep{}}
			assert.Equal(t, want, list.List)
		})
		t.Run("append no steps", func(t *testing.T) {
			t.Parallel()
			list := runstate.StepList{List: []steps.Step{}}
			list.Append()
			assert.Equal(t, []steps.Step{}, list.List)
		})
	})

	t.Run("AppendList", func(t *testing.T) {
		t.Parallel()
		t.Run("append a populated list", func(t *testing.T) {
			t.Parallel()
			list := runstate.StepList{List: []steps.Step{&steps.AbortMergeStep{}}}
			other := runstate.StepList{List: []steps.Step{&steps.StashOpenChangesStep{}}}
			list.AppendList(other)
			want := []steps.Step{&steps.AbortMergeStep{}, &steps.StashOpenChangesStep{}}
			assert.Equal(t, want, list.List)
		})
		t.Run("append an empty list", func(t *testing.T) {
			t.Parallel()
			list := runstate.StepList{List: []steps.Step{&steps.AbortMergeStep{}}}
			other := runstate.StepList{List: []steps.Step{}}
			list.AppendList(other)
			assert.Equal(t, []steps.Step{&steps.AbortMergeStep{}}, list.List)
		})
	})

	t.Run("IsEmpty", func(t *testing.T) {
		t.Parallel()
		t.Run("list is empty", func(t *testing.T) {
			t.Parallel()
			list := runstate.StepList{List: []steps.Step{}}
			assert.True(t, list.IsEmpty())
		})
		t.Run("list is not empty", func(t *testing.T) {
			t.Parallel()
			list := runstate.StepList{List: []steps.Step{&steps.AbortMergeStep{}}}
			assert.False(t, list.IsEmpty())
		})
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		list := runstate.StepList{List: []steps.Step{
			&steps.AbortMergeStep{},
			&steps.StashOpenChangesStep{},
		}}
		have, err := json.MarshalIndent(list, "", "  ")
		assert.Nil(t, err)
		// NOTE: Why does it not serialize the type names here?
		// This somehow works when serializing a StepList as part of a larger containing structure like a RunState,
		// but it doesn't work here for some reason.
		want := `
{
  "List": [
    {},
    {}
  ]
}`[1:]
		assert.Equal(t, want, string(have))
	})

	t.Run("Peek", func(t *testing.T) {
		t.Parallel()
		t.Run("populated list", func(t *testing.T) {
			t.Parallel()
			list := runstate.StepList{List: []steps.Step{&steps.AbortMergeStep{}, &steps.StashOpenChangesStep{}}}
			have := list.Peek()
			assert.Equal(t, &steps.AbortMergeStep{}, have, "returns the first element of the list")
			wantList := runstate.StepList{List: []steps.Step{&steps.AbortMergeStep{}, &steps.StashOpenChangesStep{}}}
			assert.Equal(t, wantList, list, "does not modify the list")
		})
		t.Run("empty list", func(t *testing.T) {
			t.Parallel()
			list := runstate.StepList{List: []steps.Step{}}
			have := list.Peek()
			assert.Equal(t, nil, have)
			wantList := runstate.StepList{List: []steps.Step{}}
			assert.Equal(t, wantList, list)
		})
	})

	t.Run("Pop", func(t *testing.T) {
		t.Parallel()
		t.Run("populated list", func(t *testing.T) {
			t.Parallel()
			list := runstate.StepList{List: []steps.Step{&steps.AbortMergeStep{}, &steps.StashOpenChangesStep{}}}
			have := list.Pop()
			assert.Equal(t, &steps.AbortMergeStep{}, have, "returns the first element of the list")
			wantList := runstate.StepList{List: []steps.Step{&steps.StashOpenChangesStep{}}}
			assert.Equal(t, wantList, list, "remotes the popped element from the list")
		})
		t.Run("empty list", func(t *testing.T) {
			t.Parallel()
			list := runstate.StepList{List: []steps.Step{}}
			have := list.Pop()
			assert.Equal(t, nil, have, "returns nil")
			wantList := runstate.StepList{List: []steps.Step{}}
			assert.Equal(t, wantList, list, "remotes the popped element from the list")
		})
	})

	t.Run("Prepend", func(t *testing.T) {
		t.Parallel()
		t.Run("prepend a single step", func(t *testing.T) {
			t.Parallel()
			list := runstate.StepList{List: []steps.Step{&steps.AbortMergeStep{}}}
			list.Prepend(&steps.StashOpenChangesStep{})
			want := []steps.Step{&steps.StashOpenChangesStep{}, &steps.AbortMergeStep{}}
			assert.Equal(t, want, list.List)
		})
		t.Run("prepend multiple steps", func(t *testing.T) {
			t.Parallel()
			list := runstate.StepList{List: []steps.Step{&steps.AbortMergeStep{}}}
			list.Prepend(&steps.AbortRebaseStep{}, &steps.StashOpenChangesStep{})
			want := []steps.Step{&steps.AbortRebaseStep{}, &steps.StashOpenChangesStep{}, &steps.AbortMergeStep{}}
			assert.Equal(t, want, list.List)
		})
		t.Run("prepend no steps", func(t *testing.T) {
			t.Parallel()
			list := runstate.StepList{List: []steps.Step{}}
			list.Prepend()
			assert.Equal(t, []steps.Step{}, list.List)
		})
	})

	t.Run("PrependList", func(t *testing.T) {
		t.Parallel()
		t.Run("prepend a populated list", func(t *testing.T) {
			t.Parallel()
			list := runstate.StepList{List: []steps.Step{&steps.AbortMergeStep{}}}
			other := runstate.StepList{List: []steps.Step{&steps.StashOpenChangesStep{}, &steps.RestoreOpenChangesStep{}}}
			list.PrependList(other)
			want := []steps.Step{&steps.StashOpenChangesStep{}, &steps.RestoreOpenChangesStep{}, &steps.AbortMergeStep{}}
			assert.Equal(t, want, list.List)
		})
		t.Run("prepend an empty list", func(t *testing.T) {
			t.Parallel()
			list := runstate.StepList{List: []steps.Step{&steps.AbortMergeStep{}}}
			other := runstate.StepList{List: []steps.Step{}}
			list.PrependList(other)
			want := []steps.Step{&steps.AbortMergeStep{}}
			assert.Equal(t, want, list.List)
		})
	})

	t.Run("RemoveAllButLast", func(t *testing.T) {
		t.Parallel()
		t.Run("contains the given type at the end", func(t *testing.T) {
			t.Parallel()
			have := runstate.StepList{
				List: []steps.Step{
					&steps.AbortMergeStep{},
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("branch")},
				},
			}
			have.RemoveAllButLast("*steps.CheckoutIfExistsStep")
			want := runstate.StepList{
				List: []steps.Step{
					&steps.AbortMergeStep{},
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("branch")},
				},
			}
			assert.Equal(t, want, have)
		})
		t.Run("contains the given type in the middle", func(t *testing.T) {
			t.Parallel()
			have := runstate.StepList{
				List: []steps.Step{
					&steps.AbortMergeStep{},
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("branch")},
					&steps.AbortRebaseStep{},
				},
			}
			have.RemoveAllButLast("*steps.CheckoutIfExistsStep")
			want := runstate.StepList{
				List: []steps.Step{
					&steps.AbortMergeStep{},
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("branch")},
					&steps.AbortRebaseStep{},
				},
			}
			assert.Equal(t, want, have)
		})
		t.Run("contains the given type multiple times", func(t *testing.T) {
			t.Parallel()
			have := runstate.StepList{
				List: []steps.Step{
					&steps.AbortMergeStep{},
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("branch-1")},
					&steps.AbortRebaseStep{},
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("branch-2")},
					&steps.CheckoutStep{Branch: domain.NewLocalBranchName("branch-3")},
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("branch-3")},
				},
			}
			have.RemoveAllButLast("*steps.CheckoutIfExistsStep")
			want := runstate.StepList{
				List: []steps.Step{
					&steps.AbortMergeStep{},
					&steps.AbortRebaseStep{},
					&steps.CheckoutStep{Branch: domain.NewLocalBranchName("branch-3")},
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("branch-3")},
				},
			}
			assert.Equal(t, want, have)
		})
		t.Run("does not contain the given type", func(t *testing.T) {
			t.Parallel()
			have := runstate.StepList{
				List: []steps.Step{
					&steps.AbortMergeStep{},
					&steps.AbortRebaseStep{},
					&steps.CheckoutStep{Branch: domain.NewLocalBranchName("branch-3")},
				},
			}
			have.RemoveAllButLast("*steps.CheckoutIfExistsStep")
			want := runstate.StepList{
				List: []steps.Step{
					&steps.AbortMergeStep{},
					&steps.AbortRebaseStep{},
					&steps.CheckoutStep{Branch: domain.NewLocalBranchName("branch-3")},
				},
			}
			assert.Equal(t, want, have)
		})
	})

	t.Run("RemoveDuplicateCheckoutSteps", func(t *testing.T) {
		t.Parallel()
		t.Run("has duplicate checkout steps", func(t *testing.T) {
			t.Parallel()
			give := runstate.StepList{
				List: []steps.Step{
					&steps.AbortMergeStep{},
					&steps.CheckoutStep{Branch: domain.NewLocalBranchName("branch-1")},
					&steps.CheckoutStep{Branch: domain.NewLocalBranchName("branch-2")},
				},
			}
			have := give.RemoveDuplicateCheckoutSteps()
			want := runstate.StepList{
				List: []steps.Step{
					&steps.AbortMergeStep{},
					&steps.CheckoutStep{Branch: domain.NewLocalBranchName("branch-2")},
				},
			}
			assert.Equal(t, want, have)
		})
		t.Run("has a mix of Checkout and CheckoutIfExists steps", func(t *testing.T) {
			t.Parallel()
			give := runstate.StepList{
				List: []steps.Step{
					&steps.AbortMergeStep{},
					&steps.CheckoutStep{Branch: domain.NewLocalBranchName("branch-1")},
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("branch-2")},
				},
			}
			have := give.RemoveDuplicateCheckoutSteps()
			want := runstate.StepList{
				List: []steps.Step{
					&steps.AbortMergeStep{},
					&steps.CheckoutIfExistsStep{Branch: domain.NewLocalBranchName("branch-2")},
				},
			}
			assert.Equal(t, want, have)
		})
		t.Run("has no duplicate checkout steps", func(t *testing.T) {
			t.Parallel()
			give := runstate.StepList{
				List: []steps.Step{
					&steps.AbortMergeStep{},
					&steps.AbortRebaseStep{},
				},
			}
			have := give.RemoveDuplicateCheckoutSteps()
			want := runstate.StepList{
				List: []steps.Step{
					&steps.AbortMergeStep{},
					&steps.AbortRebaseStep{},
				},
			}
			assert.Equal(t, want, have)
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		list := runstate.StepList{List: []steps.Step{
			&steps.AbortMergeStep{},
			&steps.AddToPerennialBranchesStep{
				Branch: domain.NewLocalBranchName("branch"),
			},
		}}
		have := list.String()
		want := `
StepList:
1: &steps.AbortMergeStep{EmptyStep:steps.EmptyStep{}}
2: &steps.AddToPerennialBranchesStep{Branch:domain.LocalBranchName{id:"branch"}, EmptyStep:steps.EmptyStep{}}
`[1:]
		assert.Equal(t, want, have)
	})

	t.Run("StepTypes", func(t *testing.T) {
		t.Parallel()
		list := runstate.StepList{
			List: []steps.Step{
				&steps.AbortMergeStep{},
				&steps.CheckoutStep{Branch: domain.NewLocalBranchName("branch")},
			},
		}
		have := list.StepTypes()
		want := []string{"*steps.AbortMergeStep", "*steps.CheckoutStep"}
		assert.Equal(t, want, have)
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Parallel()
		give := `
[
	{
		"data": {
			"Hard": false,
			"MustHaveSHA": "abcdef",
			"SetToSHA": "123456"
		},
		"type": "ResetCurrentBranchToSHAStep"
	},
	{
		"data": {},
		"type": "StashOpenChangesStep"
	}
]`[1:]
		have := runstate.StepList{}
		err := json.Unmarshal([]byte(give), &have)
		assert.Nil(t, err)
		want := runstate.StepList{List: []steps.Step{
			&steps.ResetCurrentBranchToSHAStep{
				Hard:        false,
				MustHaveSHA: domain.NewSHA("abcdef"),
				SetToSHA:    domain.NewSHA("123456"),
			},
			&steps.StashOpenChangesStep{},
		}}
		assert.Equal(t, want, have)
	})
}
