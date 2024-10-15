package program_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
	"github.com/git-town/git-town/v16/internal/vm/program"
	"github.com/git-town/git-town/v16/internal/vm/shared"
	"github.com/shoenig/test/must"
)

func TestProgram(t *testing.T) {
	t.Parallel()

	t.Run("Append", func(t *testing.T) {
		t.Parallel()
		t.Run("append a single opcode", func(t *testing.T) {
			t.Parallel()
			have := program.Program{&opcodes.MergeAbort{}}
			have.Add(&opcodes.StashOpenChanges{})
			want := []shared.Opcode{&opcodes.MergeAbort{}, &opcodes.StashOpenChanges{}}
			must.Eq(t, want, have)
		})
		t.Run("append multiple opcodes", func(t *testing.T) {
			t.Parallel()
			have := program.Program{&opcodes.MergeAbort{}}
			have.Add(&opcodes.RebaseAbort{}, &opcodes.StashOpenChanges{})
			want := []shared.Opcode{&opcodes.MergeAbort{}, &opcodes.RebaseAbort{}, &opcodes.StashOpenChanges{}}
			must.Eq(t, want, have)
		})
		t.Run("append no opcodes", func(t *testing.T) {
			t.Parallel()
			have := program.Program{}
			have.Add()
			must.Len(t, 0, have)
		})
	})

	t.Run("AddProgram", func(t *testing.T) {
		t.Parallel()
		t.Run("append a populated list", func(t *testing.T) {
			t.Parallel()
			have := program.Program{&opcodes.MergeAbort{}}
			other := program.Program{&opcodes.StashOpenChanges{}}
			have.AddProgram(other)
			want := []shared.Opcode{&opcodes.MergeAbort{}, &opcodes.StashOpenChanges{}}
			must.Eq(t, want, have)
		})
		t.Run("append an empty list", func(t *testing.T) {
			t.Parallel()
			have := program.Program{&opcodes.MergeAbort{}}
			other := program.Program{}
			have.AddProgram(other)
			must.Eq(t, []shared.Opcode{&opcodes.MergeAbort{}}, have)
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
			have := program.Program{&opcodes.MergeAbort{}}
			must.False(t, have.IsEmpty())
		})
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		give := program.Program{
			&opcodes.MergeAbort{},
			&opcodes.StashOpenChanges{},
		}
		have, err := json.MarshalIndent(give, "", "  ")
		must.NoError(t, err)
		want := `
[
  {
    "data": {},
    "type": "MergeAbort"
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
				&opcodes.MergeAbort{},
				&opcodes.StashOpenChanges{},
			}
			have := give.Peek()
			must.Eq(t, "*opcodes.MergeAbort", reflect.TypeOf(have).String())
			wantProgram := program.Program{&opcodes.MergeAbort{}, &opcodes.StashOpenChanges{}}
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
			give := program.Program{&opcodes.MergeAbort{}, &opcodes.StashOpenChanges{}}
			have := give.Pop()
			must.EqOp(t, "*opcodes.MergeAbort", reflect.TypeOf(have).String())
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
			give := program.Program{&opcodes.MergeAbort{}}
			give.Prepend(&opcodes.StashOpenChanges{})
			want := []shared.Opcode{&opcodes.StashOpenChanges{}, &opcodes.MergeAbort{}}
			must.Eq(t, want, give)
		})
		t.Run("prepend multiple opcodes", func(t *testing.T) {
			t.Parallel()
			give := program.Program{&opcodes.MergeAbort{}}
			give.Prepend(&opcodes.RebaseAbort{}, &opcodes.StashOpenChanges{})
			want := []shared.Opcode{&opcodes.RebaseAbort{}, &opcodes.StashOpenChanges{}, &opcodes.MergeAbort{}}
			must.Eq(t, want, give)
		})
		t.Run("prepend no opcodes", func(t *testing.T) {
			t.Parallel()
			give := program.Program{}
			give.Prepend()
			must.Len(t, 0, give)
		})
		t.Run("used as a higher-level function", func(t *testing.T) {
			t.Parallel()
			give := program.Program{&opcodes.MergeAbort{}}
			prepend := give.Prepend
			prepend(&opcodes.RebaseAbort{}, &opcodes.StashOpenChanges{})
			want := []shared.Opcode{&opcodes.RebaseAbort{}, &opcodes.StashOpenChanges{}, &opcodes.MergeAbort{}}
			must.Eq(t, want, give)
		})
	})

	t.Run("PrependProgram", func(t *testing.T) {
		t.Parallel()
		t.Run("prepend a populated list", func(t *testing.T) {
			t.Parallel()
			give := program.Program{&opcodes.MergeAbort{}}
			other := program.Program{&opcodes.StashOpenChanges{}, &opcodes.RestoreOpenChanges{}}
			give.PrependProgram(other)
			want := []shared.Opcode{&opcodes.StashOpenChanges{}, &opcodes.RestoreOpenChanges{}, &opcodes.MergeAbort{}}
			must.Eq(t, want, give)
		})
		t.Run("prepend an empty list", func(t *testing.T) {
			t.Parallel()
			give := program.Program{&opcodes.MergeAbort{}}
			other := program.Program{}
			give.PrependProgram(other)
			want := []shared.Opcode{&opcodes.MergeAbort{}}
			must.Eq(t, want, give)
		})
	})

	t.Run("RemoveAllButLast", func(t *testing.T) {
		t.Parallel()
		t.Run("contains the given type at the end", func(t *testing.T) {
			t.Parallel()
			give := program.Program{
				&opcodes.MergeAbort{},
				&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch")},
			}
			have := give.RemoveAllButLast("*opcodes.CheckoutIfExists")
			want := program.Program{
				&opcodes.MergeAbort{},
				&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch")},
			}
			must.Eq(t, want, have)
		})
		t.Run("contains the given type in the middle", func(t *testing.T) {
			t.Parallel()
			give := program.Program{
				&opcodes.MergeAbort{},
				&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.RebaseAbort{},
			}
			have := give.RemoveAllButLast("*opcodes.CheckoutIfExists")
			want := program.Program{
				&opcodes.MergeAbort{},
				&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch")},
				&opcodes.RebaseAbort{},
			}
			must.Eq(t, want, have)
		})
		t.Run("contains the given type multiple times", func(t *testing.T) {
			t.Parallel()
			give := program.Program{
				&opcodes.MergeAbort{},
				&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch-1")},
				&opcodes.RebaseAbort{},
				&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch-2")},
				&opcodes.CheckoutIfNeeded{Branch: gitdomain.NewLocalBranchName("branch-3")},
				&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch-3")},
			}
			have := give.RemoveAllButLast("*opcodes.CheckoutIfExists")
			want := program.Program{
				&opcodes.MergeAbort{},
				&opcodes.RebaseAbort{},
				&opcodes.CheckoutIfNeeded{Branch: gitdomain.NewLocalBranchName("branch-3")},
				&opcodes.CheckoutIfExists{Branch: gitdomain.NewLocalBranchName("branch-3")},
			}
			must.Eq(t, want, have)
		})
		t.Run("does not contain the given type", func(t *testing.T) {
			t.Parallel()
			give := program.Program{
				&opcodes.MergeAbort{},
				&opcodes.RebaseAbort{},
				&opcodes.CheckoutIfNeeded{Branch: gitdomain.NewLocalBranchName("branch-3")},
			}
			have := give.RemoveAllButLast("*opcodes.CheckoutIfExists")
			want := program.Program{
				&opcodes.MergeAbort{},
				&opcodes.RebaseAbort{},
				&opcodes.CheckoutIfNeeded{Branch: gitdomain.NewLocalBranchName("branch-3")},
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		give := program.Program{
			&opcodes.MergeAbort{},
			&opcodes.BranchesPerennialAdd{
				Branch: gitdomain.NewLocalBranchName("branch"),
			},
		}
		have := give.String()
		want := `
Program:
1: &opcodes.MergeAbort{undeclaredOpcodeMethods:opcodes.undeclaredOpcodeMethods{}}
2: &opcodes.BranchesPerennialAdd{Branch:"branch", undeclaredOpcodeMethods:opcodes.undeclaredOpcodeMethods{}}
`[1:]
		must.EqOp(t, want, have)
	})

	t.Run("OpcodeTypes", func(t *testing.T) {
		t.Parallel()
		prog := program.Program{
			&opcodes.MergeAbort{},
			&opcodes.CheckoutIfNeeded{Branch: gitdomain.NewLocalBranchName("branch")},
		}
		have := prog.OpcodeTypes()
		want := []string{"*opcodes.MergeAbort", "*opcodes.CheckoutIfNeeded"}
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
		"type": "BranchCurrentResetToSHAIfNeeded"
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
			&opcodes.BranchCurrentResetToSHAIfNeeded{
				Hard:        false,
				MustHaveSHA: gitdomain.NewSHA("abcdef"),
				SetToSHA:    gitdomain.NewSHA("123456"),
			},
			&opcodes.StashOpenChanges{},
		}
		must.Eq(t, want, have)
	})
}
