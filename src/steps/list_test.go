package steps_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/step"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/shoenig/test/must"
)

func TestList(t *testing.T) {
	t.Parallel()

	t.Run("Append", func(t *testing.T) {
		t.Parallel()
		t.Run("append a single step", func(t *testing.T) {
			t.Parallel()
			list := steps.List{List: []step.Step{&step.AbortMerge{}}}
			list.Add(&step.StashOpenChanges{})
			want := []step.Step{&step.AbortMerge{}, &step.StashOpenChanges{}}
			must.Eq(t, want, list.List)
		})
		t.Run("append multiple steps", func(t *testing.T) {
			t.Parallel()
			list := steps.List{List: []step.Step{&step.AbortMerge{}}}
			list.Add(&step.AbortRebase{}, &step.StashOpenChanges{})
			want := []step.Step{&step.AbortMerge{}, &step.AbortRebase{}, &step.StashOpenChanges{}}
			must.Eq(t, want, list.List)
		})
		t.Run("append no steps", func(t *testing.T) {
			t.Parallel()
			list := steps.List{List: []step.Step{}}
			list.Add()
			must.Eq(t, []step.Step{}, list.List)
		})
	})

	t.Run("AppendList", func(t *testing.T) {
		t.Parallel()
		t.Run("append a populated list", func(t *testing.T) {
			t.Parallel()
			list := steps.List{List: []step.Step{&step.AbortMerge{}}}
			other := steps.List{List: []step.Step{&step.StashOpenChanges{}}}
			list.AddList(other)
			want := []step.Step{&step.AbortMerge{}, &step.StashOpenChanges{}}
			must.Eq(t, want, list.List)
		})
		t.Run("append an empty list", func(t *testing.T) {
			t.Parallel()
			list := steps.List{List: []step.Step{&step.AbortMerge{}}}
			other := steps.List{List: []step.Step{}}
			list.AddList(other)
			must.Eq(t, []step.Step{&step.AbortMerge{}}, list.List)
		})
	})

	t.Run("IsEmpty", func(t *testing.T) {
		t.Parallel()
		t.Run("list is empty", func(t *testing.T) {
			t.Parallel()
			list := steps.List{List: []step.Step{}}
			must.True(t, list.IsEmpty())
		})
		t.Run("list is not empty", func(t *testing.T) {
			t.Parallel()
			list := steps.List{List: []step.Step{&step.AbortMerge{}}}
			must.False(t, list.IsEmpty())
		})
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		list := steps.List{List: []step.Step{
			&step.AbortMerge{},
			&step.StashOpenChanges{},
		}}
		have, err := json.MarshalIndent(list, "", "  ")
		must.NoError(t, err)
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
		must.EqOp(t, want, string(have))
	})

	t.Run("Peek", func(t *testing.T) {
		t.Parallel()
		t.Run("populated list", func(t *testing.T) {
			t.Parallel()
			list := steps.List{List: []step.Step{&step.AbortMerge{}, &step.StashOpenChanges{}}}
			have := list.Peek()
			must.Eq(t, "*step.AbortMerge", reflect.TypeOf(have).String())
			wantList := steps.List{List: []step.Step{&step.AbortMerge{}, &step.StashOpenChanges{}}}
			must.Eq(t, wantList, list)
		})
		t.Run("empty list", func(t *testing.T) {
			t.Parallel()
			list := steps.List{List: []step.Step{}}
			have := list.Peek()
			must.EqOp(t, nil, have)
			wantList := steps.List{List: []step.Step{}}
			must.Eq(t, wantList, list)
		})
	})

	t.Run("Pop", func(t *testing.T) {
		t.Parallel()
		t.Run("populated list", func(t *testing.T) {
			t.Parallel()
			list := steps.List{List: []step.Step{&step.AbortMerge{}, &step.StashOpenChanges{}}}
			have := list.Pop()
			must.EqOp(t, "*step.AbortMerge", reflect.TypeOf(have).String())
			wantList := steps.List{List: []step.Step{&step.StashOpenChanges{}}}
			must.Eq(t, wantList, list)
		})
		t.Run("empty list", func(t *testing.T) {
			t.Parallel()
			list := steps.List{List: []step.Step{}}
			have := list.Pop()
			must.EqOp(t, nil, have)
			wantList := steps.List{List: []step.Step{}}
			must.Eq(t, wantList, list)
		})
	})

	t.Run("Prepend", func(t *testing.T) {
		t.Parallel()
		t.Run("prepend a single step", func(t *testing.T) {
			t.Parallel()
			list := steps.List{List: []step.Step{&step.AbortMerge{}}}
			list.Prepend(&step.StashOpenChanges{})
			want := []step.Step{&step.StashOpenChanges{}, &step.AbortMerge{}}
			must.Eq(t, want, list.List)
		})
		t.Run("prepend multiple steps", func(t *testing.T) {
			t.Parallel()
			list := steps.List{List: []step.Step{&step.AbortMerge{}}}
			list.Prepend(&step.AbortRebase{}, &step.StashOpenChanges{})
			want := []step.Step{&step.AbortRebase{}, &step.StashOpenChanges{}, &step.AbortMerge{}}
			must.Eq(t, want, list.List)
		})
		t.Run("prepend no steps", func(t *testing.T) {
			t.Parallel()
			list := steps.List{List: []step.Step{}}
			list.Prepend()
			must.Eq(t, []step.Step{}, list.List)
		})
		t.Run("used as callback", func(t *testing.T) {
			t.Parallel()
			list := steps.List{List: []step.Step{&step.AbortMerge{}}}
			prepend := list.Prepend
			prepend(&step.AbortRebase{}, &step.StashOpenChanges{})
			want := []step.Step{&step.AbortRebase{}, &step.StashOpenChanges{}, &step.AbortMerge{}}
			must.Eq(t, want, list.List)
		})
	})

	t.Run("PrependList", func(t *testing.T) {
		t.Parallel()
		t.Run("prepend a populated list", func(t *testing.T) {
			t.Parallel()
			list := steps.List{List: []step.Step{&step.AbortMerge{}}}
			other := steps.List{List: []step.Step{&step.StashOpenChanges{}, &step.RestoreOpenChanges{}}}
			list.PrependList(other)
			want := []step.Step{&step.StashOpenChanges{}, &step.RestoreOpenChanges{}, &step.AbortMerge{}}
			must.Eq(t, want, list.List)
		})
		t.Run("prepend an empty list", func(t *testing.T) {
			t.Parallel()
			list := steps.List{List: []step.Step{&step.AbortMerge{}}}
			other := steps.List{List: []step.Step{}}
			list.PrependList(other)
			want := []step.Step{&step.AbortMerge{}}
			must.Eq(t, want, list.List)
		})
	})

	t.Run("RemoveAllButLast", func(t *testing.T) {
		t.Parallel()
		t.Run("contains the given type at the end", func(t *testing.T) {
			t.Parallel()
			have := steps.List{
				List: []step.Step{
					&step.AbortMerge{},
					&step.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch")},
				},
			}
			have.RemoveAllButLast("*step.CheckoutIfExists")
			want := steps.List{
				List: []step.Step{
					&step.AbortMerge{},
					&step.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch")},
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("contains the given type in the middle", func(t *testing.T) {
			t.Parallel()
			have := steps.List{
				List: []step.Step{
					&step.AbortMerge{},
					&step.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch")},
					&step.AbortRebase{},
				},
			}
			have.RemoveAllButLast("*step.CheckoutIfExists")
			want := steps.List{
				List: []step.Step{
					&step.AbortMerge{},
					&step.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch")},
					&step.AbortRebase{},
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("contains the given type multiple times", func(t *testing.T) {
			t.Parallel()
			have := steps.List{
				List: []step.Step{
					&step.AbortMerge{},
					&step.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch-1")},
					&step.AbortRebase{},
					&step.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch-2")},
					&step.Checkout{Branch: domain.NewLocalBranchName("branch-3")},
					&step.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch-3")},
				},
			}
			have.RemoveAllButLast("*step.CheckoutIfExists")
			want := steps.List{
				List: []step.Step{
					&step.AbortMerge{},
					&step.AbortRebase{},
					&step.Checkout{Branch: domain.NewLocalBranchName("branch-3")},
					&step.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch-3")},
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("does not contain the given type", func(t *testing.T) {
			t.Parallel()
			have := steps.List{
				List: []step.Step{
					&step.AbortMerge{},
					&step.AbortRebase{},
					&step.Checkout{Branch: domain.NewLocalBranchName("branch-3")},
				},
			}
			have.RemoveAllButLast("*step.CheckoutIfExists")
			want := steps.List{
				List: []step.Step{
					&step.AbortMerge{},
					&step.AbortRebase{},
					&step.Checkout{Branch: domain.NewLocalBranchName("branch-3")},
				},
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("RemoveDuplicateCheckoutSteps", func(t *testing.T) {
		t.Parallel()
		t.Run("has duplicate checkout steps", func(t *testing.T) {
			t.Parallel()
			give := steps.List{
				List: []step.Step{
					&step.AbortMerge{},
					&step.Checkout{Branch: domain.NewLocalBranchName("branch-1")},
					&step.Checkout{Branch: domain.NewLocalBranchName("branch-2")},
				},
			}
			have := give.RemoveDuplicateCheckoutSteps()
			want := steps.List{
				List: []step.Step{
					&step.AbortMerge{},
					&step.Checkout{Branch: domain.NewLocalBranchName("branch-2")},
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("has a mix of Checkout and CheckoutIfExists steps", func(t *testing.T) {
			t.Parallel()
			give := steps.List{
				List: []step.Step{
					&step.AbortMerge{},
					&step.Checkout{Branch: domain.NewLocalBranchName("branch-1")},
					&step.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch-2")},
				},
			}
			have := give.RemoveDuplicateCheckoutSteps()
			want := steps.List{
				List: []step.Step{
					&step.AbortMerge{},
					&step.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch-2")},
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("has no duplicate checkout steps", func(t *testing.T) {
			t.Parallel()
			give := steps.List{
				List: []step.Step{
					&step.AbortMerge{},
					&step.AbortRebase{},
				},
			}
			have := give.RemoveDuplicateCheckoutSteps()
			want := steps.List{
				List: []step.Step{
					&step.AbortMerge{},
					&step.AbortRebase{},
				},
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		list := steps.List{List: []step.Step{
			&step.AbortMerge{},
			&step.AddToPerennialBranches{
				Branch: domain.NewLocalBranchName("branch"),
			},
		}}
		have := list.String()
		want := `
StepList:
1: &step.AbortMerge{Empty:step.Empty{}}
2: &step.AddToPerennialBranches{Branch:domain.LocalBranchName{id:"branch"}, Empty:step.Empty{}}
`[1:]
		must.EqOp(t, want, have)
	})

	t.Run("StepTypes", func(t *testing.T) {
		t.Parallel()
		list := steps.List{
			List: []step.Step{
				&step.AbortMerge{},
				&step.Checkout{Branch: domain.NewLocalBranchName("branch")},
			},
		}
		have := list.StepTypes()
		want := []string{"*step.AbortMerge", "*step.Checkout"}
		must.Eq(t, want, have)
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
		"type": "ResetCurrentBranchToSHA"
	},
	{
		"data": {},
		"type": "StashOpenChanges"
	}
]`[1:]
		have := steps.List{}
		err := json.Unmarshal([]byte(give), &have)
		must.NoError(t, err)
		want := steps.List{List: []step.Step{
			&step.ResetCurrentBranchToSHA{
				Hard:        false,
				MustHaveSHA: domain.NewSHA("abcdef"),
				SetToSHA:    domain.NewSHA("123456"),
			},
			&step.StashOpenChanges{},
		}}
		must.Eq(t, want, have)
	})
}
