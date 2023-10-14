package program_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/program"
	"github.com/git-town/git-town/v9/src/vm/step"
	"github.com/shoenig/test/must"
)

func TestProgram(t *testing.T) {
	t.Parallel()

	t.Run("Append", func(t *testing.T) {
		t.Parallel()
		t.Run("append a single step", func(t *testing.T) {
			t.Parallel()
			have := program.Program{Steps: []step.Step{&step.AbortMerge{}}}
			have.Add(&step.StashOpenChanges{})
			want := []step.Step{&step.AbortMerge{}, &step.StashOpenChanges{}}
			must.Eq(t, want, have.Steps)
		})
		t.Run("append multiple steps", func(t *testing.T) {
			t.Parallel()
			have := program.Program{Steps: []step.Step{&step.AbortMerge{}}}
			have.Add(&step.AbortRebase{}, &step.StashOpenChanges{})
			want := []step.Step{&step.AbortMerge{}, &step.AbortRebase{}, &step.StashOpenChanges{}}
			must.Eq(t, want, have.Steps)
		})
		t.Run("append no steps", func(t *testing.T) {
			t.Parallel()
			have := program.Program{Steps: []step.Step{}}
			have.Add()
			must.Eq(t, []step.Step{}, have.Steps)
		})
	})

	t.Run("AddProgram", func(t *testing.T) {
		t.Parallel()
		t.Run("append a populated list", func(t *testing.T) {
			t.Parallel()
			have := program.Program{Steps: []step.Step{&step.AbortMerge{}}}
			other := program.Program{Steps: []step.Step{&step.StashOpenChanges{}}}
			have.AddProgram(other)
			want := []step.Step{&step.AbortMerge{}, &step.StashOpenChanges{}}
			must.Eq(t, want, have.Steps)
		})
		t.Run("append an empty list", func(t *testing.T) {
			t.Parallel()
			have := program.Program{Steps: []step.Step{&step.AbortMerge{}}}
			other := program.Program{Steps: []step.Step{}}
			have.AddProgram(other)
			must.Eq(t, []step.Step{&step.AbortMerge{}}, have.Steps)
		})
	})

	t.Run("IsEmpty", func(t *testing.T) {
		t.Parallel()
		t.Run("list is empty", func(t *testing.T) {
			t.Parallel()
			have := program.Program{Steps: []step.Step{}}
			must.True(t, have.IsEmpty())
		})
		t.Run("list is not empty", func(t *testing.T) {
			t.Parallel()
			have := program.Program{Steps: []step.Step{&step.AbortMerge{}}}
			must.False(t, have.IsEmpty())
		})
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		give := program.Program{Steps: []step.Step{
			&step.AbortMerge{},
			&step.StashOpenChanges{},
		}}
		have, err := json.MarshalIndent(give, "", "  ")
		must.NoError(t, err)
		// NOTE: Why does it not serialize the type names here?
		// This somehow works when serializing a program as part of a larger containing structure like a RunState,
		// but it doesn't work here for some reason.
		want := `
{
  "Steps": [
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
			give := program.Program{Steps: []step.Step{&step.AbortMerge{}, &step.StashOpenChanges{}}}
			have := give.Peek()
			must.Eq(t, "*step.AbortMerge", reflect.TypeOf(have).String())
			wantProgram := program.Program{Steps: []step.Step{&step.AbortMerge{}, &step.StashOpenChanges{}}}
			must.Eq(t, wantProgram, give)
		})
		t.Run("empty list", func(t *testing.T) {
			t.Parallel()
			give := program.Program{Steps: []step.Step{}}
			have := give.Peek()
			must.EqOp(t, nil, have)
			wantProgram := program.Program{Steps: []step.Step{}}
			must.Eq(t, wantProgram, give)
		})
	})

	t.Run("Pop", func(t *testing.T) {
		t.Parallel()
		t.Run("populated list", func(t *testing.T) {
			t.Parallel()
			give := program.Program{Steps: []step.Step{&step.AbortMerge{}, &step.StashOpenChanges{}}}
			have := give.Pop()
			must.EqOp(t, "*step.AbortMerge", reflect.TypeOf(have).String())
			wantProgram := program.Program{Steps: []step.Step{&step.StashOpenChanges{}}}
			must.Eq(t, wantProgram, give)
		})
		t.Run("empty list", func(t *testing.T) {
			t.Parallel()
			give := program.Program{Steps: []step.Step{}}
			have := give.Pop()
			must.EqOp(t, nil, have)
			wantProgram := program.Program{Steps: []step.Step{}}
			must.Eq(t, wantProgram, give)
		})
	})

	t.Run("Prepend", func(t *testing.T) {
		t.Parallel()
		t.Run("prepend a single step", func(t *testing.T) {
			t.Parallel()
			give := program.Program{Steps: []step.Step{&step.AbortMerge{}}}
			give.Prepend(&step.StashOpenChanges{})
			want := []step.Step{&step.StashOpenChanges{}, &step.AbortMerge{}}
			must.Eq(t, want, give.Steps)
		})
		t.Run("prepend multiple steps", func(t *testing.T) {
			t.Parallel()
			give := program.Program{Steps: []step.Step{&step.AbortMerge{}}}
			give.Prepend(&step.AbortRebase{}, &step.StashOpenChanges{})
			want := []step.Step{&step.AbortRebase{}, &step.StashOpenChanges{}, &step.AbortMerge{}}
			must.Eq(t, want, give.Steps)
		})
		t.Run("prepend no steps", func(t *testing.T) {
			t.Parallel()
			give := program.Program{Steps: []step.Step{}}
			give.Prepend()
			must.Eq(t, []step.Step{}, give.Steps)
		})
		t.Run("used as a higher-level function", func(t *testing.T) {
			t.Parallel()
			give := program.Program{Steps: []step.Step{&step.AbortMerge{}}}
			prepend := give.Prepend
			prepend(&step.AbortRebase{}, &step.StashOpenChanges{})
			want := []step.Step{&step.AbortRebase{}, &step.StashOpenChanges{}, &step.AbortMerge{}}
			must.Eq(t, want, give.Steps)
		})
	})

	t.Run("PrependProgram", func(t *testing.T) {
		t.Parallel()
		t.Run("prepend a populated list", func(t *testing.T) {
			t.Parallel()
			give := program.Program{Steps: []step.Step{&step.AbortMerge{}}}
			other := program.Program{Steps: []step.Step{&step.StashOpenChanges{}, &step.RestoreOpenChanges{}}}
			give.PrependProgram(other)
			want := []step.Step{&step.StashOpenChanges{}, &step.RestoreOpenChanges{}, &step.AbortMerge{}}
			must.Eq(t, want, give.Steps)
		})
		t.Run("prepend an empty list", func(t *testing.T) {
			t.Parallel()
			give := program.Program{Steps: []step.Step{&step.AbortMerge{}}}
			other := program.Program{Steps: []step.Step{}}
			give.PrependProgram(other)
			want := []step.Step{&step.AbortMerge{}}
			must.Eq(t, want, give.Steps)
		})
	})

	t.Run("RemoveAllButLast", func(t *testing.T) {
		t.Parallel()
		t.Run("contains the given type at the end", func(t *testing.T) {
			t.Parallel()
			have := program.Program{
				Steps: []step.Step{
					&step.AbortMerge{},
					&step.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch")},
				},
			}
			have.RemoveAllButLast("*step.CheckoutIfExists")
			want := program.Program{
				Steps: []step.Step{
					&step.AbortMerge{},
					&step.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch")},
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("contains the given type in the middle", func(t *testing.T) {
			t.Parallel()
			have := program.Program{
				Steps: []step.Step{
					&step.AbortMerge{},
					&step.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch")},
					&step.AbortRebase{},
				},
			}
			have.RemoveAllButLast("*step.CheckoutIfExists")
			want := program.Program{
				Steps: []step.Step{
					&step.AbortMerge{},
					&step.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch")},
					&step.AbortRebase{},
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("contains the given type multiple times", func(t *testing.T) {
			t.Parallel()
			have := program.Program{
				Steps: []step.Step{
					&step.AbortMerge{},
					&step.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch-1")},
					&step.AbortRebase{},
					&step.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch-2")},
					&step.Checkout{Branch: domain.NewLocalBranchName("branch-3")},
					&step.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch-3")},
				},
			}
			have.RemoveAllButLast("*step.CheckoutIfExists")
			want := program.Program{
				Steps: []step.Step{
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
			have := program.Program{
				Steps: []step.Step{
					&step.AbortMerge{},
					&step.AbortRebase{},
					&step.Checkout{Branch: domain.NewLocalBranchName("branch-3")},
				},
			}
			have.RemoveAllButLast("*step.CheckoutIfExists")
			want := program.Program{
				Steps: []step.Step{
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
			give := program.Program{
				Steps: []step.Step{
					&step.AbortMerge{},
					&step.Checkout{Branch: domain.NewLocalBranchName("branch-1")},
					&step.Checkout{Branch: domain.NewLocalBranchName("branch-2")},
				},
			}
			have := give.RemoveDuplicateCheckoutSteps()
			want := program.Program{
				Steps: []step.Step{
					&step.AbortMerge{},
					&step.Checkout{Branch: domain.NewLocalBranchName("branch-2")},
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("has a mix of Checkout and CheckoutIfExists steps", func(t *testing.T) {
			t.Parallel()
			give := program.Program{
				Steps: []step.Step{
					&step.AbortMerge{},
					&step.Checkout{Branch: domain.NewLocalBranchName("branch-1")},
					&step.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch-2")},
				},
			}
			have := give.RemoveDuplicateCheckoutSteps()
			want := program.Program{
				Steps: []step.Step{
					&step.AbortMerge{},
					&step.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch-2")},
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("has no duplicate checkout steps", func(t *testing.T) {
			t.Parallel()
			give := program.Program{
				Steps: []step.Step{
					&step.AbortMerge{},
					&step.AbortRebase{},
				},
			}
			have := give.RemoveDuplicateCheckoutSteps()
			want := program.Program{
				Steps: []step.Step{
					&step.AbortMerge{},
					&step.AbortRebase{},
				},
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		give := program.Program{Steps: []step.Step{
			&step.AbortMerge{},
			&step.AddToPerennialBranches{
				Branch: domain.NewLocalBranchName("branch"),
			},
		}}
		have := give.String()
		want := `
Program:
1: &step.AbortMerge{Empty:step.Empty{}}
2: &step.AddToPerennialBranches{Branch:domain.LocalBranchName{id:"branch"}, Empty:step.Empty{}}
`[1:]
		must.EqOp(t, want, have)
	})

	t.Run("StepTypes", func(t *testing.T) {
		t.Parallel()
		list := program.Program{
			Steps: []step.Step{
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
		have := program.Program{}
		err := json.Unmarshal([]byte(give), &have)
		must.NoError(t, err)
		want := program.Program{Steps: []step.Step{
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
