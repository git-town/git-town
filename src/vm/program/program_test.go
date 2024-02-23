package program_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/opcodes"
	"github.com/git-town/git-town/v12/src/vm/program"
	"github.com/git-town/git-town/v12/src/vm/shared"
	"github.com/shoenig/test/must"
)

func TestProgram(t *testing.T) {
	t.Parallel()

	t.Run("Append", func(t *testing.T) {
		t.Parallel()
		t.Run("append a single opcode", func(t *testing.T) {
			t.Parallel()
			have := program.Program{&opcodes.AbortMerge{}}
			have.Add(&opcodes.StashOpenChanges{})
			want := []shared.Opcode{&opcodes.AbortMerge{}, &opcodes.StashOpenChanges{}}
			must.Eq(t, want, have)
		})
		t.Run("append multiple opcodes", func(t *testing.T) {
			t.Parallel()
			have := program.Program{&opcodes.AbortMerge{}}
			have.Add(&opcodes.AbortRebase{}, &opcodes.StashOpenChanges{})
			want := []shared.Opcode{&opcodes.AbortMerge{}, &opcodes.AbortRebase{}, &opcodes.StashOpenChanges{}}
			must.Eq(t, want, have)
		})
		t.Run("append no opcodes", func(t *testing.T) {
			t.Parallel()
			have := program.Program{}
			have.Add()
			must.Eq(t, []shared.Opcode{}, have)
		})
	})

	t.Run("AddProgram", func(t *testing.T) {
		t.Parallel()
		t.Run("append a populated list", func(t *testing.T) {
			t.Parallel()
			have := program.Program{&opcodes.AbortMerge{}}
			other := program.Program{&opcodes.StashOpenChanges{}}
			have.AddProgram(other)
			want := []shared.Opcode{&opcodes.AbortMerge{}, &opcodes.StashOpenChanges{}}
			must.Eq(t, want, have)
		})
		t.Run("append an empty list", func(t *testing.T) {
			t.Parallel()
			have := program.Program{&opcodes.AbortMerge{}}
			other := program.Program{}
			have.AddProgram(other)
			must.Eq(t, []shared.Opcode{&opcodes.AbortMerge{}}, have)
		})
	})

	t.Run("IsEmpty", func(t *testing.T) {
		t.Parallel()
		t.Run("list is empty", func(t *testing.T) {
			t.Parallel()
			have := program.Program{}
			must.True(t, have.IsEmpty())
		})
		t.Run("list is not empty", func(t *testing.T) {
			t.Parallel()
			have := program.Program{&opcodes.AbortMerge{}}
			must.False(t, have.IsEmpty())
		})
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		give := program.Program{
			&opcodes.AbortMerge{},
			&opcodes.StashOpenChanges{},
		}
		have, err := json.MarshalIndent(give, "", "  ")
		must.NoError(t, err)
		want := `
[
  {
    "data": {},
    "type": "AbortMerge"
  },
  {
    "data": {},
    "type": "StashOpenChanges"
  }
]`[1:]
		must.EqOp(t, want, string(have))
	})

	t.Run("Peek", func(t *testing.T) {
		t.Parallel()
		t.Run("populated list", func(t *testing.T) {
			t.Parallel()
			give := program.Program{
				&opcodes.AbortMerge{},
				&opcodes.StashOpenChanges{},
			}
			have := give.Peek()
			must.Eq(t, "*opcodes.AbortMerge", reflect.TypeOf(have).String())
			wantProgram := program.Program{&opcodes.AbortMerge{}, &opcodes.StashOpenChanges{}}
			must.Eq(t, wantProgram, give)
		})
		t.Run("empty list", func(t *testing.T) {
			t.Parallel()
			give := program.Program{}
			have := give.Peek()
			must.EqOp(t, nil, have)
			wantProgram := program.Program{}
			must.Eq(t, wantProgram, give)
		})
	})

	t.Run("Pop", func(t *testing.T) {
		t.Parallel()
		t.Run("populated list", func(t *testing.T) {
			t.Parallel()
			give := program.Program{&opcodes.AbortMerge{}, &opcodes.StashOpenChanges{}}
			have := give.Pop()
			must.EqOp(t, "*opcodes.AbortMerge", reflect.TypeOf(have).String())
			wantProgram := program.Program{&opcodes.StashOpenChanges{}}
			must.Eq(t, wantProgram, give)
		})
		t.Run("empty list", func(t *testing.T) {
			t.Parallel()
			give := program.Program{}
			have := give.Pop()
			must.EqOp(t, nil, have)
			wantProgram := program.Program{}
			must.Eq(t, wantProgram, give)
		})
	})

	t.Run("Prepend", func(t *testing.T) {
		t.Parallel()
		t.Run("prepend a single opcode", func(t *testing.T) {
			t.Parallel()
			give := program.Program{&opcodes.AbortMerge{}}
			give.Prepend(&opcodes.StashOpenChanges{})
			want := []shared.Opcode{&opcodes.StashOpenChanges{}, &opcodes.AbortMerge{}}
			must.Eq(t, want, give)
		})
		t.Run("prepend multiple opcodes", func(t *testing.T) {
			t.Parallel()
			give := program.Program{&opcodes.AbortMerge{}}
			give.Prepend(&opcodes.AbortRebase{}, &opcodes.StashOpenChanges{})
			want := []shared.Opcode{&opcodes.AbortRebase{}, &opcodes.StashOpenChanges{}, &opcodes.AbortMerge{}}
			must.Eq(t, want, give)
		})
		t.Run("prepend no opcodes", func(t *testing.T) {
			t.Parallel()
			give := program.Program{}
			give.Prepend()
			must.Eq(t, []shared.Opcode{}, give)
		})
		t.Run("used as a higher-level function", func(t *testing.T) {
			t.Parallel()
			give := program.Program{&opcodes.AbortMerge{}}
			prepend := give.Prepend
			prepend(&opcodes.AbortRebase{}, &opcodes.StashOpenChanges{})
			want := []shared.Opcode{&opcodes.AbortRebase{}, &opcodes.StashOpenChanges{}, &opcodes.AbortMerge{}}
			must.Eq(t, want, give)
		})
	})

	t.Run("PrependProgram", func(t *testing.T) {
		t.Parallel()
		t.Run("prepend a populated list", func(t *testing.T) {
			t.Parallel()
			give := program.Program{&opcodes.AbortMerge{}}
			other := program.Program{&opcodes.StashOpenChanges{}, &opcodes.RestoreOpenChanges{}}
			give.PrependProgram(other)
			want := []shared.Opcode{&opcodes.StashOpenChanges{}, &opcodes.RestoreOpenChanges{}, &opcodes.AbortMerge{}}
			must.Eq(t, want, give)
		})
		t.Run("prepend an empty list", func(t *testing.T) {
			t.Parallel()
			give := program.Program{&opcodes.AbortMerge{}}
			other := program.Program{}
			give.PrependProgram(other)
			want := []shared.Opcode{&opcodes.AbortMerge{}}
			must.Eq(t, want, give)
		})
	})

	t.Run("RemoveAllButLast", func(t *testing.T) {
		t.Parallel()
		t.Run("contains the given type at the end", func(t *testing.T) {
			t.Parallel()
			give := program.Program{
				&opcodes.AbortMerge{},
				&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch")},
			}
			have := give.RemoveAllButLast("*opcodes.CheckoutIfExists")
			want := program.Program{
				&opcodes.AbortMerge{},
				&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch")},
			}
			must.Eq(t, want, have)
		})
		t.Run("contains the given type in the middle", func(t *testing.T) {
			t.Parallel()
			give := program.Program{
				&opcodes.AbortMerge{},
				&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.AbortRebase{},
			}
			have := give.RemoveAllButLast("*opcodes.CheckoutIfExists")
			want := program.Program{
				&opcodes.AbortMerge{},
				&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.AbortRebase{},
			}
			must.Eq(t, want, have)
		})
		t.Run("contains the given type multiple times", func(t *testing.T) {
			t.Parallel()
			give := program.Program{
				&opcodes.AbortMerge{},
				&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch-1")},
				&opcodes.AbortRebase{},
				&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch-2")},
				&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("branch-3")},
				&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch-3")},
			}
			have := give.RemoveAllButLast("*opcodes.CheckoutIfExists")
			want := program.Program{
				&opcodes.AbortMerge{},
				&opcodes.AbortRebase{},
				&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("branch-3")},
				&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch-3")},
			}
			must.Eq(t, want, have)
		})
		t.Run("does not contain the given type", func(t *testing.T) {
			t.Parallel()
			give := program.Program{
				&opcodes.AbortMerge{},
				&opcodes.AbortRebase{},
				&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("branch-3")},
			}
			have := give.RemoveAllButLast("*opcodes.CheckoutIfExists")
			want := program.Program{
				&opcodes.AbortMerge{},
				&opcodes.AbortRebase{},
				&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("branch-3")},
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("RemoveDuplicateCheckout", func(t *testing.T) {
		t.Parallel()
		t.Run("has duplicate checkout opcodes", func(t *testing.T) {
			t.Parallel()
			give := program.Program{
				&opcodes.AbortMerge{},
				&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("branch-1")},
				&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("branch-2")},
				&opcodes.AbortRebase{},
			}
			give.RemoveDuplicateCheckout()
			want := program.Program{
				&opcodes.AbortMerge{},
				&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("branch-2")},
				&opcodes.AbortRebase{},
			}
			must.Eq(t, want, give)
		})
		t.Run("has duplicate checkout opcodes mixed with end-of-branch opcodes", func(t *testing.T) {
			t.Parallel()
			give := program.Program{
				&opcodes.AbortMerge{},
				&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("branch-1")},
				&opcodes.EndOfBranchProgram{},
				&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("branch-2")},
				&opcodes.EndOfBranchProgram{},
				&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("branch-3")},
				&opcodes.AbortRebase{},
			}
			give.RemoveDuplicateCheckout()
			want := program.Program{
				&opcodes.AbortMerge{},
				&opcodes.EndOfBranchProgram{},
				&opcodes.EndOfBranchProgram{},
				&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("branch-3")},
				&opcodes.AbortRebase{},
			}
			must.Eq(t, want, give)
		})
		t.Run("has a mix of Checkout and CheckoutIfExists opcodes", func(t *testing.T) {
			t.Parallel()
			give := program.Program{
				&opcodes.AbortMerge{},
				&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("branch-1")},
				&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch-2")},
			}
			give.RemoveDuplicateCheckout()
			want := program.Program{
				&opcodes.AbortMerge{},
				&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch-2")},
			}
			must.Eq(t, want, give)
		})
		t.Run("has no duplicate checkout opcodes", func(t *testing.T) {
			t.Parallel()
			give := program.Program{
				&opcodes.AbortMerge{},
				&opcodes.AbortRebase{},
			}
			give.RemoveDuplicateCheckout()
			want := program.Program{
				&opcodes.AbortMerge{},
				&opcodes.AbortRebase{},
			}
			must.Eq(t, want, give)
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		give := program.Program{
			&opcodes.AbortMerge{},
			&opcodes.AddToPerennialBranches{
				Branch: gitdomain.NewLocalBranchName("branch"),
			},
		}
		have := give.String()
		want := `
Program:
1: &opcodes.AbortMerge{undeclaredOpcodeMethods:opcodes.undeclaredOpcodeMethods{}}
2: &opcodes.AddToPerennialBranches{Branch:"branch", undeclaredOpcodeMethods:opcodes.undeclaredOpcodeMethods{}}
`[1:]
		must.EqOp(t, want, have)
	})

	t.Run("OpcodeTypes", func(t *testing.T) {
		t.Parallel()
		prog := program.Program{
			&opcodes.AbortMerge{},
			&opcodes.Checkout{Branch: gitdomain.NewLocalBranchName("branch")},
		}
		have := prog.OpcodeTypes()
		want := []string{"*opcodes.AbortMerge", "*opcodes.Checkout"}
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
		want := program.Program{
			&opcodes.ResetCurrentBranchToSHA{
				Hard:        false,
				MustHaveSHA: gitdomain.NewSHA("abcdef"),
				SetToSHA:    gitdomain.NewSHA("123456"),
			},
			&opcodes.StashOpenChanges{},
		}
		must.Eq(t, want, have)
	})
}
