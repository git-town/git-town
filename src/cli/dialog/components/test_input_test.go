package components_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/shoenig/test/must"
)

func TestTestInputs(t *testing.T) {
	t.Parallel()

	t.Run("LoadTestInputs", func(t *testing.T) {
		t.Parallel()
		give := []string{
			"foo=bar",
			"GITTOWN_DIALOG_INPUT_1=enter",
			"GITTOWN_DIALOG_INPUT_2=space|down|space|5|enter",
			"GITTOWN_DIALOG_INPUT_3=ctrl+c",
		}
		have := components.LoadTestInputs(give)
		want := components.NewTestInputs(
			components.TestInput{tea.KeyMsg{Type: tea.KeyEnter}},
			components.TestInput{
				tea.KeyMsg{Type: tea.KeySpace},
				tea.KeyMsg{Type: tea.KeyDown},
				tea.KeyMsg{Type: tea.KeySpace},
				tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}},
				tea.KeyMsg{Type: tea.KeyEnter},
			},
			components.TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}},
		)
		must.Eq(t, want, have)
	})

	t.Run("TestInputs.Next", func(t *testing.T) {
		t.Parallel()
		testInputs := components.NewTestInputs(
			components.TestInput{tea.KeyMsg{Type: tea.KeyCtrlA}},
			components.TestInput{tea.KeyMsg{Type: tea.KeyCtrlB}},
			components.TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}},
		)
		// request the first entry: A
		haveNext := testInputs.Next()
		wantNext := components.TestInput{tea.KeyMsg{Type: tea.KeyCtrlA}}
		must.Eq(t, wantNext, haveNext)
		wantRemaining := components.NewTestInputs(
			components.TestInput{tea.KeyMsg{Type: tea.KeyCtrlB}},
			components.TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}},
		)
		must.Eq(t, wantRemaining, testInputs)
		// request the next entry: B
		haveNext = testInputs.Next()
		wantNext = components.TestInput{tea.KeyMsg{Type: tea.KeyCtrlB}}
		must.Eq(t, wantNext, haveNext)
		wantRemaining = components.NewTestInputs(
			components.TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}},
		)
		must.Eq(t, wantRemaining, testInputs)
		// request the next entry: C
		haveNext = testInputs.Next()
		wantNext = components.TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}}
		must.Eq(t, wantNext, haveNext)
		must.EqOp(t, 0, testInputs.Len())
		// request the next entry: empty
		haveNext = testInputs.Next()
		wantNext = components.TestInput{}
		must.Eq(t, wantNext, haveNext)
		must.EqOp(t, 0, testInputs.Len())
	})

	t.Run("ParseTestInput", func(t *testing.T) {
		t.Parallel()
		t.Run("multiple values", func(t *testing.T) {
			t.Parallel()
			have := components.ParseTestInput("enter|space|ctrl+c")
			want := components.TestInput{
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
			have := components.ParseTestInput("enter")
			want := components.TestInput{
				tea.KeyMsg{
					Type: tea.KeyEnter,
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			have := components.ParseTestInput("")
			want := components.TestInput{}
			must.Eq(t, want, have)
		})
	})
}
