package program_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/vm/opcode"
	"github.com/git-town/git-town/v11/src/vm/program"
	"github.com/git-town/git-town/v11/src/vm/shared"
	"github.com/shoenig/test/must"
)

func TestProgram(t *testing.T) {
	t.Parallel()

	t.Run("Append", func(t *testing.T) {
		t.Parallel()
		t.Run("append a single opcode", func(t *testing.T) {
			t.Parallel()
			have := program.Program{&opcode.AbortMerge{}}
			have.Add(&opcode.StashOpenChanges{})
			want := []shared.Opcode{&opcode.AbortMerge{}, &opcode.StashOpenChanges{}}
			must.Eq(t, want, have)
		})
		t.Run("append multiple opcodes", func(t *testing.T) {
			t.Parallel()
			have := program.Program{&opcode.AbortMerge{}}
			have.Add(&opcode.AbortRebase{}, &opcode.StashOpenChanges{})
			want := []shared.Opcode{&opcode.AbortMerge{}, &opcode.AbortRebase{}, &opcode.StashOpenChanges{}}
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
			have := program.Program{&opcode.AbortMerge{}}
			other := program.Program{&opcode.StashOpenChanges{}}
			have.AddProgram(other)
			want := []shared.Opcode{&opcode.AbortMerge{}, &opcode.StashOpenChanges{}}
			must.Eq(t, want, have)
		})
		t.Run("append an empty list", func(t *testing.T) {
			t.Parallel()
			have := program.Program{&opcode.AbortMerge{}}
			other := program.Program{}
			have.AddProgram(other)
			must.Eq(t, []shared.Opcode{&opcode.AbortMerge{}}, have)
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
			have := program.Program{&opcode.AbortMerge{}}
			must.False(t, have.IsEmpty())
		})
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		give := program.Program{
			&opcode.AbortMerge{},
			&opcode.StashOpenChanges{},
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

	t.Run("MoveToEnd", func(t *testing.T) {
		t.Parallel()

		t.Run("program contains opcode at the last position", func(t *testing.T) {
			t.Parallel()
			have := program.Program{
				&opcode.AbortMerge{},
				&opcode.AbortRebase{},
			}
			have.MoveToEnd(&opcode.AbortRebase{})
			want := program.Program{
				&opcode.AbortMerge{},
				&opcode.AbortRebase{},
			}
			must.Eq(t, want, have)
		})

		t.Run("program contains element in the middle", func(t *testing.T) {
			t.Parallel()
			have := program.Program{
				&opcode.AbortMerge{},
				&opcode.RestoreOpenChanges{},
				&opcode.AbortRebase{},
			}
			have.MoveToEnd(&opcode.RestoreOpenChanges{})
			want := []shared.Opcode{
				&opcode.AbortMerge{},
				&opcode.AbortRebase{},
				&opcode.RestoreOpenChanges{},
			}
			must.Eq(t, want, have)
		})

		t.Run("program does not contain the element", func(t *testing.T) {
			t.Parallel()
			have := program.Program{
				&opcode.AbortMerge{},
			}
			have.MoveToEnd(&opcode.ContinueMerge{})
			want := program.Program{
				&opcode.AbortMerge{},
			}
			must.Eq(t, want, have)
		})

		t.Run("multiple occurrences of the opcode to move", func(t *testing.T) {
			t.Parallel()
			have := program.Program{
				&opcode.ContinueMerge{},
				&opcode.AbortMerge{},
				&opcode.ContinueMerge{},
				&opcode.AbortRebase{},
				&opcode.ContinueMerge{},
			}
			have.MoveToEnd(&opcode.ContinueMerge{})
			want := program.Program{
				&opcode.AbortMerge{},
				&opcode.AbortRebase{},
				&opcode.ContinueMerge{},
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("Peek", func(t *testing.T) {
		t.Parallel()
		t.Run("populated list", func(t *testing.T) {
			t.Parallel()
			give := program.Program{
				&opcode.AbortMerge{},
				&opcode.StashOpenChanges{},
			}
			have := give.Peek()
			must.Eq(t, "*opcode.AbortMerge", reflect.TypeOf(have).String())
			wantProgram := program.Program{&opcode.AbortMerge{}, &opcode.StashOpenChanges{}}
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
			give := program.Program{&opcode.AbortMerge{}, &opcode.StashOpenChanges{}}
			have := give.Pop()
			must.EqOp(t, "*opcode.AbortMerge", reflect.TypeOf(have).String())
			wantProgram := program.Program{&opcode.StashOpenChanges{}}
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
			give := program.Program{&opcode.AbortMerge{}}
			give.Prepend(&opcode.StashOpenChanges{})
			want := []shared.Opcode{&opcode.StashOpenChanges{}, &opcode.AbortMerge{}}
			must.Eq(t, want, give)
		})
		t.Run("prepend multiple opcodes", func(t *testing.T) {
			t.Parallel()
			give := program.Program{&opcode.AbortMerge{}}
			give.Prepend(&opcode.AbortRebase{}, &opcode.StashOpenChanges{})
			want := []shared.Opcode{&opcode.AbortRebase{}, &opcode.StashOpenChanges{}, &opcode.AbortMerge{}}
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
			give := program.Program{&opcode.AbortMerge{}}
			prepend := give.Prepend
			prepend(&opcode.AbortRebase{}, &opcode.StashOpenChanges{})
			want := []shared.Opcode{&opcode.AbortRebase{}, &opcode.StashOpenChanges{}, &opcode.AbortMerge{}}
			must.Eq(t, want, give)
		})
	})

	t.Run("PrependProgram", func(t *testing.T) {
		t.Parallel()
		t.Run("prepend a populated list", func(t *testing.T) {
			t.Parallel()
			give := program.Program{&opcode.AbortMerge{}}
			other := program.Program{&opcode.StashOpenChanges{}, &opcode.RestoreOpenChanges{}}
			give.PrependProgram(other)
			want := []shared.Opcode{&opcode.StashOpenChanges{}, &opcode.RestoreOpenChanges{}, &opcode.AbortMerge{}}
			must.Eq(t, want, give)
		})
		t.Run("prepend an empty list", func(t *testing.T) {
			t.Parallel()
			give := program.Program{&opcode.AbortMerge{}}
			other := program.Program{}
			give.PrependProgram(other)
			want := []shared.Opcode{&opcode.AbortMerge{}}
			must.Eq(t, want, give)
		})
	})

	t.Run("RemoveAllButLast", func(t *testing.T) {
		t.Parallel()
		t.Run("contains the given type at the end", func(t *testing.T) {
			t.Parallel()
			have := program.Program{
				&opcode.AbortMerge{},
				&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch")},
			}
			have.RemoveAllButLast("*opcode.CheckoutIfExists")
			want := program.Program{
				&opcode.AbortMerge{},
				&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch")},
			}
			must.Eq(t, want, have)
		})
		t.Run("contains the given type in the middle", func(t *testing.T) {
			t.Parallel()
			have := program.Program{
				&opcode.AbortMerge{},
				&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch")},
				&opcode.AbortRebase{},
			}
			have.RemoveAllButLast("*opcode.CheckoutIfExists")
			want := program.Program{
				&opcode.AbortMerge{},
				&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch")},
				&opcode.AbortRebase{},
			}
			must.Eq(t, want, have)
		})
		t.Run("contains the given type multiple times", func(t *testing.T) {
			t.Parallel()
			have := program.Program{
				&opcode.AbortMerge{},
				&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch-1")},
				&opcode.AbortRebase{},
				&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch-2")},
				&opcode.Checkout{Branch: domain.NewLocalBranchName("branch-3")},
				&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch-3")},
			}
			have.RemoveAllButLast("*opcode.CheckoutIfExists")
			want := program.Program{
				&opcode.AbortMerge{},
				&opcode.AbortRebase{},
				&opcode.Checkout{Branch: domain.NewLocalBranchName("branch-3")},
				&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch-3")},
			}
			must.Eq(t, want, have)
		})
		t.Run("does not contain the given type", func(t *testing.T) {
			t.Parallel()
			have := program.Program{
				&opcode.AbortMerge{},
				&opcode.AbortRebase{},
				&opcode.Checkout{Branch: domain.NewLocalBranchName("branch-3")},
			}
			have.RemoveAllButLast("*opcode.CheckoutIfExists")
			want := program.Program{
				&opcode.AbortMerge{},
				&opcode.AbortRebase{},
				&opcode.Checkout{Branch: domain.NewLocalBranchName("branch-3")},
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("RemoveDuplicateCheckout", func(t *testing.T) {
		t.Parallel()
		t.Run("has duplicate checkout opcodes", func(t *testing.T) {
			t.Parallel()
			give := program.Program{
				&opcode.AbortMerge{},
				&opcode.Checkout{Branch: domain.NewLocalBranchName("branch-1")},
				&opcode.Checkout{Branch: domain.NewLocalBranchName("branch-2")},
			}
			give.RemoveDuplicateCheckout()
			want := program.Program{
				&opcode.AbortMerge{},
				&opcode.Checkout{Branch: domain.NewLocalBranchName("branch-2")},
			}
			must.Eq(t, want, give)
		})
		t.Run("has a mix of Checkout and CheckoutIfExists opcodes", func(t *testing.T) {
			t.Parallel()
			give := program.Program{
				&opcode.AbortMerge{},
				&opcode.Checkout{Branch: domain.NewLocalBranchName("branch-1")},
				&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch-2")},
			}
			give.RemoveDuplicateCheckout()
			want := program.Program{
				&opcode.AbortMerge{},
				&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch-2")},
			}
			must.Eq(t, want, give)
		})
		t.Run("has no duplicate checkout opcodes", func(t *testing.T) {
			t.Parallel()
			give := program.Program{
				&opcode.AbortMerge{},
				&opcode.AbortRebase{},
			}
			give.RemoveDuplicateCheckout()
			want := program.Program{
				&opcode.AbortMerge{},
				&opcode.AbortRebase{},
			}
			must.Eq(t, want, give)
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		give := program.Program{
			&opcode.AbortMerge{},
			&opcode.AddToPerennialBranches{
				Branch: domain.NewLocalBranchName("branch"),
			},
		}
		have := give.String()
		want := `
Program:
1: &opcode.AbortMerge{undeclaredOpcodeMethods:opcode.undeclaredOpcodeMethods{}}
2: &opcode.AddToPerennialBranches{Branch:"branch", undeclaredOpcodeMethods:opcode.undeclaredOpcodeMethods{}}
`[1:]
		must.EqOp(t, want, have)
	})

	t.Run("OpcodeTypes", func(t *testing.T) {
		t.Parallel()
		prog := program.Program{
			&opcode.AbortMerge{},
			&opcode.Checkout{Branch: domain.NewLocalBranchName("branch")},
		}
		have := prog.OpcodeTypes()
		want := []string{"*opcode.AbortMerge", "*opcode.Checkout"}
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
			&opcode.ResetCurrentBranchToSHA{
				Hard:        false,
				MustHaveSHA: domain.NewSHA("abcdef"),
				SetToSHA:    domain.NewSHA("123456"),
			},
			&opcode.StashOpenChanges{},
		}
		must.Eq(t, want, have)
	})
}
