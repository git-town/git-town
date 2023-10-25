package opcode_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/opcode"
	"github.com/git-town/git-town/v9/src/vm/shared"
	"github.com/shoenig/test/must"
)

func TestJSON(t *testing.T) {
	t.Parallel()

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		jsonstep := opcode.JSON{
			Opcode: &opcode.Checkout{
				Branch: domain.NewLocalBranchName("branch-1"),
			},
		}
		have, err := json.MarshalIndent(jsonstep, "", "  ")
		must.NoError(t, err)
		// NOTE: It's unclear why this doesn't contain the "data" and "type" fields from JSONStep's MarshalJSON method here.
		//       Marshaling an entire RunState somehow works correctly.
		want := `
{
  "Opcode": {
    "Branch": "branch-1"
  }
}`[1:]
		must.EqOp(t, want, string(have))
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Parallel()
		give := `
{
	"data": {
    "Branch": "branch-1"
  },
	"type": "Checkout"
}`[1:]
		have := opcode.JSON{
			Opcode: &opcode.Checkout{
				Branch: domain.EmptyLocalBranchName(),
			},
		}
		err := json.Unmarshal([]byte(give), &have)
		must.NoError(t, err)
		want := opcode.JSON{
			Opcode: &opcode.Checkout{
				Branch: domain.NewLocalBranchName("branch-1"),
			},
		}
		must.Eq(t, want, have)
	})
}

func TestLookup(t *testing.T) {
	t.Parallel()
	tests := map[string]shared.Opcode{
		"AbortMerge":  &opcode.AbortMerge{},
		"AbortRebase": &opcode.AbortRebase{},
		"unknown":     nil,
	}
	for give, want := range tests {
		have := opcode.Lookup(give)
		must.EqOp(t, want, have)
	}
}

func TestProgram(t *testing.T) {
	t.Parallel()

	t.Run("Append", func(t *testing.T) {
		t.Parallel()
		t.Run("append a single opcode", func(t *testing.T) {
			t.Parallel()
			have := opcode.Program{Opcodes: []shared.Opcode{&opcode.AbortMerge{}}}
			have.Add(&opcode.StashOpenChanges{})
			want := []shared.Opcode{&opcode.AbortMerge{}, &opcode.StashOpenChanges{}}
			must.Eq(t, want, have.Opcodes)
		})
		t.Run("append multiple opcodes", func(t *testing.T) {
			t.Parallel()
			have := opcode.Program{Opcodes: []shared.Opcode{&opcode.AbortMerge{}}}
			have.Add(&opcode.AbortRebase{}, &opcode.StashOpenChanges{})
			want := []shared.Opcode{&opcode.AbortMerge{}, &opcode.AbortRebase{}, &opcode.StashOpenChanges{}}
			must.Eq(t, want, have.Opcodes)
		})
		t.Run("append no opcodes", func(t *testing.T) {
			t.Parallel()
			have := opcode.Program{Opcodes: []shared.Opcode{}}
			have.Add()
			must.Eq(t, []shared.Opcode{}, have.Opcodes)
		})
	})

	t.Run("AddProgram", func(t *testing.T) {
		t.Parallel()
		t.Run("append a populated list", func(t *testing.T) {
			t.Parallel()
			have := opcode.Program{Opcodes: []shared.Opcode{&opcode.AbortMerge{}}}
			other := opcode.Program{Opcodes: []shared.Opcode{&opcode.StashOpenChanges{}}}
			have.AddProgram(other)
			want := []shared.Opcode{&opcode.AbortMerge{}, &opcode.StashOpenChanges{}}
			must.Eq(t, want, have.Opcodes)
		})
		t.Run("append an empty list", func(t *testing.T) {
			t.Parallel()
			have := opcode.Program{Opcodes: []shared.Opcode{&opcode.AbortMerge{}}}
			other := opcode.Program{Opcodes: []shared.Opcode{}}
			have.AddProgram(other)
			must.Eq(t, []shared.Opcode{&opcode.AbortMerge{}}, have.Opcodes)
		})
	})

	t.Run("IsEmpty", func(t *testing.T) {
		t.Parallel()
		t.Run("list is empty", func(t *testing.T) {
			t.Parallel()
			have := opcode.Program{Opcodes: []shared.Opcode{}}
			must.True(t, have.IsEmpty())
		})
		t.Run("list is not empty", func(t *testing.T) {
			t.Parallel()
			have := opcode.Program{Opcodes: []shared.Opcode{&opcode.AbortMerge{}}}
			must.False(t, have.IsEmpty())
		})
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		give := opcode.Program{Opcodes: []shared.Opcode{
			&opcode.AbortMerge{},
			&opcode.StashOpenChanges{},
		}}
		have, err := json.MarshalIndent(give, "", "  ")
		must.NoError(t, err)
		// NOTE: Why does it not serialize the type names here?
		// This somehow works when serializing a program as part of a larger containing structure like a RunState,
		// but it doesn't work here for some reason.
		want := `
{
  "Opcodes": [
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
			give := opcode.Program{Opcodes: []shared.Opcode{&opcode.AbortMerge{}, &opcode.StashOpenChanges{}}}
			have := give.Peek()
			must.Eq(t, "*opcode.AbortMerge", reflect.TypeOf(have).String())
			wantProgram := opcode.Program{Opcodes: []shared.Opcode{&opcode.AbortMerge{}, &opcode.StashOpenChanges{}}}
			must.Eq(t, wantProgram, give)
		})
		t.Run("empty list", func(t *testing.T) {
			t.Parallel()
			give := opcode.Program{Opcodes: []shared.Opcode{}}
			have := give.Peek()
			must.EqOp(t, nil, have)
			wantProgram := opcode.Program{Opcodes: []shared.Opcode{}}
			must.Eq(t, wantProgram, give)
		})
	})

	t.Run("Pop", func(t *testing.T) {
		t.Parallel()
		t.Run("populated list", func(t *testing.T) {
			t.Parallel()
			give := opcode.Program{Opcodes: []shared.Opcode{&opcode.AbortMerge{}, &opcode.StashOpenChanges{}}}
			have := give.Pop()
			must.EqOp(t, "*opcode.AbortMerge", reflect.TypeOf(have).String())
			wantProgram := opcode.Program{Opcodes: []shared.Opcode{&opcode.StashOpenChanges{}}}
			must.Eq(t, wantProgram, give)
		})
		t.Run("empty list", func(t *testing.T) {
			t.Parallel()
			give := opcode.Program{Opcodes: []shared.Opcode{}}
			have := give.Pop()
			must.EqOp(t, nil, have)
			wantProgram := opcode.Program{Opcodes: []shared.Opcode{}}
			must.Eq(t, wantProgram, give)
		})
	})

	t.Run("Prepend", func(t *testing.T) {
		t.Parallel()
		t.Run("prepend a single opcode", func(t *testing.T) {
			t.Parallel()
			give := opcode.Program{Opcodes: []shared.Opcode{&opcode.AbortMerge{}}}
			give.Prepend(&opcode.StashOpenChanges{})
			want := []shared.Opcode{&opcode.StashOpenChanges{}, &opcode.AbortMerge{}}
			must.Eq(t, want, give.Opcodes)
		})
		t.Run("prepend multiple opcodes", func(t *testing.T) {
			t.Parallel()
			give := opcode.Program{Opcodes: []shared.Opcode{&opcode.AbortMerge{}}}
			give.Prepend(&opcode.AbortRebase{}, &opcode.StashOpenChanges{})
			want := []shared.Opcode{&opcode.AbortRebase{}, &opcode.StashOpenChanges{}, &opcode.AbortMerge{}}
			must.Eq(t, want, give.Opcodes)
		})
		t.Run("prepend no opcodes", func(t *testing.T) {
			t.Parallel()
			give := opcode.Program{Opcodes: []shared.Opcode{}}
			give.Prepend()
			must.Eq(t, []shared.Opcode{}, give.Opcodes)
		})
		t.Run("used as a higher-level function", func(t *testing.T) {
			t.Parallel()
			give := opcode.Program{Opcodes: []shared.Opcode{&opcode.AbortMerge{}}}
			prepend := give.Prepend
			prepend(&opcode.AbortRebase{}, &opcode.StashOpenChanges{})
			want := []shared.Opcode{&opcode.AbortRebase{}, &opcode.StashOpenChanges{}, &opcode.AbortMerge{}}
			must.Eq(t, want, give.Opcodes)
		})
	})

	t.Run("PrependProgram", func(t *testing.T) {
		t.Parallel()
		t.Run("prepend a populated list", func(t *testing.T) {
			t.Parallel()
			give := opcode.Program{Opcodes: []shared.Opcode{&opcode.AbortMerge{}}}
			other := opcode.Program{Opcodes: []shared.Opcode{&opcode.StashOpenChanges{}, &opcode.RestoreOpenChanges{}}}
			give.PrependProgram(other)
			want := []shared.Opcode{&opcode.StashOpenChanges{}, &opcode.RestoreOpenChanges{}, &opcode.AbortMerge{}}
			must.Eq(t, want, give.Opcodes)
		})
		t.Run("prepend an empty list", func(t *testing.T) {
			t.Parallel()
			give := opcode.Program{Opcodes: []shared.Opcode{&opcode.AbortMerge{}}}
			other := opcode.Program{Opcodes: []shared.Opcode{}}
			give.PrependProgram(other)
			want := []shared.Opcode{&opcode.AbortMerge{}}
			must.Eq(t, want, give.Opcodes)
		})
	})

	t.Run("RemoveAllButLast", func(t *testing.T) {
		t.Parallel()
		t.Run("contains the given type at the end", func(t *testing.T) {
			t.Parallel()
			have := opcode.Program{
				Opcodes: []shared.Opcode{
					&opcode.AbortMerge{},
					&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch")},
				},
			}
			have.RemoveAllButLast("*opcode.CheckoutIfExists")
			want := opcode.Program{
				Opcodes: []shared.Opcode{
					&opcode.AbortMerge{},
					&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch")},
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("contains the given type in the middle", func(t *testing.T) {
			t.Parallel()
			have := opcode.Program{
				Opcodes: []shared.Opcode{
					&opcode.AbortMerge{},
					&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch")},
					&opcode.AbortRebase{},
				},
			}
			have.RemoveAllButLast("*opcode.CheckoutIfExists")
			want := opcode.Program{
				Opcodes: []shared.Opcode{
					&opcode.AbortMerge{},
					&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch")},
					&opcode.AbortRebase{},
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("contains the given type multiple times", func(t *testing.T) {
			t.Parallel()
			have := opcode.Program{
				Opcodes: []shared.Opcode{
					&opcode.AbortMerge{},
					&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch-1")},
					&opcode.AbortRebase{},
					&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch-2")},
					&opcode.Checkout{Branch: domain.NewLocalBranchName("branch-3")},
					&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch-3")},
				},
			}
			have.RemoveAllButLast("*opcode.CheckoutIfExists")
			want := opcode.Program{
				Opcodes: []shared.Opcode{
					&opcode.AbortMerge{},
					&opcode.AbortRebase{},
					&opcode.Checkout{Branch: domain.NewLocalBranchName("branch-3")},
					&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch-3")},
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("does not contain the given type", func(t *testing.T) {
			t.Parallel()
			have := opcode.Program{
				Opcodes: []shared.Opcode{
					&opcode.AbortMerge{},
					&opcode.AbortRebase{},
					&opcode.Checkout{Branch: domain.NewLocalBranchName("branch-3")},
				},
			}
			have.RemoveAllButLast("*opcode.CheckoutIfExists")
			want := opcode.Program{
				Opcodes: []shared.Opcode{
					&opcode.AbortMerge{},
					&opcode.AbortRebase{},
					&opcode.Checkout{Branch: domain.NewLocalBranchName("branch-3")},
				},
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("RemoveDuplicateCheckout", func(t *testing.T) {
		t.Parallel()
		t.Run("has duplicate checkout opcodes", func(t *testing.T) {
			t.Parallel()
			give := opcode.Program{
				Opcodes: []shared.Opcode{
					&opcode.AbortMerge{},
					&opcode.Checkout{Branch: domain.NewLocalBranchName("branch-1")},
					&opcode.Checkout{Branch: domain.NewLocalBranchName("branch-2")},
				},
			}
			have := give.RemoveDuplicateCheckout()
			want := opcode.Program{
				Opcodes: []shared.Opcode{
					&opcode.AbortMerge{},
					&opcode.Checkout{Branch: domain.NewLocalBranchName("branch-2")},
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("has a mix of Checkout and CheckoutIfExists opcodes", func(t *testing.T) {
			t.Parallel()
			give := opcode.Program{
				Opcodes: []shared.Opcode{
					&opcode.AbortMerge{},
					&opcode.Checkout{Branch: domain.NewLocalBranchName("branch-1")},
					&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch-2")},
				},
			}
			have := give.RemoveDuplicateCheckout()
			want := opcode.Program{
				Opcodes: []shared.Opcode{
					&opcode.AbortMerge{},
					&opcode.CheckoutIfExists{Branch: domain.NewLocalBranchName("branch-2")},
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("has no duplicate checkout opcodes", func(t *testing.T) {
			t.Parallel()
			give := opcode.Program{
				Opcodes: []shared.Opcode{
					&opcode.AbortMerge{},
					&opcode.AbortRebase{},
				},
			}
			have := give.RemoveDuplicateCheckout()
			want := opcode.Program{
				Opcodes: []shared.Opcode{
					&opcode.AbortMerge{},
					&opcode.AbortRebase{},
				},
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		give := opcode.Program{Opcodes: []shared.Opcode{
			&opcode.AbortMerge{},
			&opcode.AddToPerennialBranches{
				Branch: domain.NewLocalBranchName("branch"),
			},
		}}
		have := give.String()
		want := `
Program:
1: &opcode.AbortMerge{undeclaredOpcodeMethods:opcode.undeclaredOpcodeMethods{}}
2: &opcode.AddToPerennialBranches{Branch:domain.LocalBranchName{id:"branch"}, undeclaredOpcodeMethods:opcode.undeclaredOpcodeMethods{}}
`[1:]
		must.EqOp(t, want, have)
	})

	t.Run("OpcodeTypes", func(t *testing.T) {
		t.Parallel()
		prog := opcode.Program{
			Opcodes: []shared.Opcode{
				&opcode.AbortMerge{},
				&opcode.Checkout{Branch: domain.NewLocalBranchName("branch")},
			},
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
		have := opcode.Program{}
		err := json.Unmarshal([]byte(give), &have)
		must.NoError(t, err)
		want := opcode.Program{Opcodes: []shared.Opcode{
			&opcode.ResetCurrentBranchToSHA{
				Hard:        false,
				MustHaveSHA: domain.NewSHA("abcdef"),
				SetToSHA:    domain.NewSHA("123456"),
			},
			&opcode.StashOpenChanges{},
		}}
		must.Eq(t, want, have)
	})
}
