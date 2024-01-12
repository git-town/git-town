package dialog_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/shoenig/test/must"
)

func TestTestInputs(t *testing.T) {
	t.Parallel()

	t.Run("LoadTestInputs", func(t *testing.T) {
		t.Parallel()
		give := []string{
			"foo=bar",
			"GITTOWN_DIALOG_INPUT1=enter",
			"GITTOWN_DIALOG_INPUT2=space|down|space|5|enter",
			"GITTOWN_DIALOG_INPUT2=ctrl-c",
		}
		have := dialog.LoadTestInputs(give)
		want := dialog.TestInputs{
			dialog.TestInput{tea.KeyMsg{Type: tea.KeyEnter}},
			dialog.TestInput{
				tea.KeyMsg{Type: tea.KeySpace},
				tea.KeyMsg{Type: tea.KeyDown},
				tea.KeyMsg{Type: tea.KeySpace},
				tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}},
				tea.KeyMsg{Type: tea.KeyEnter},
			},
			dialog.TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}},
		}
		must.Eq(t, want, have)
	})

	t.Run("TestInputs.Next", func(t *testing.T) {
		t.Parallel()
		testInputs := dialog.TestInputs{
			dialog.TestInput{tea.KeyMsg{Type: tea.KeyCtrlA}},
			dialog.TestInput{tea.KeyMsg{Type: tea.KeyCtrlB}},
			dialog.TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}},
		}
		haveNext := testInputs.Next()
		wantNext := dialog.TestInput{tea.KeyMsg{Type: tea.KeyCtrlA}}
		must.Eq(t, wantNext, haveNext)
		wantRemaining := dialog.TestInputs{
			dialog.TestInput{tea.KeyMsg{Type: tea.KeyCtrlB}},
			dialog.TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}},
		}
		must.Eq(t, wantRemaining, testInputs)
	})

	t.Run("ParseTestInput", func(t *testing.T) {
		t.Parallel()
		t.Run("multiple values", func(t *testing.T) {
			t.Parallel()
			have := dialog.ParseTestInput("enter|space|ctrl-c")
			want := dialog.TestInput{
				tea.KeyMsg{
					Type: tea.KeyEnter,
				},
				tea.KeyMsg{
					Type: tea.KeySpace,
				},
				tea.KeyMsg{
					Type: tea.KeyCtrlC,
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("single value", func(t *testing.T) {
			t.Parallel()
			have := dialog.ParseTestInput("enter")
			want := dialog.TestInput{
				tea.KeyMsg{
					Type: tea.KeyEnter,
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			have := dialog.ParseTestInput("enter")
			want := dialog.TestInput{}
			must.Eq(t, want, have)
		})
	})
}
