package dialogcomponents_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
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
		have := dialogcomponents.LoadTestInputs(give)
		want := dialogcomponents.NewTestInputs(
			dialogcomponents.TestInput{tea.KeyMsg{Type: tea.KeyEnter}},
			dialogcomponents.TestInput{
				tea.KeyMsg{Type: tea.KeySpace},
				tea.KeyMsg{Type: tea.KeyDown},
				tea.KeyMsg{Type: tea.KeySpace},
				tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}},
				tea.KeyMsg{Type: tea.KeyEnter},
			},
			dialogcomponents.TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}},
		)
		must.Eq(t, want, have)
	})

	t.Run("ParseTestInput", func(t *testing.T) {
		t.Parallel()
		t.Run("multiple values", func(t *testing.T) {
			t.Parallel()
			have := dialogcomponents.ParseTestInput("enter|space|ctrl+c")
			want := dialogcomponents.TestInput{
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
			have := dialogcomponents.ParseTestInput("enter")
			want := dialogcomponents.TestInput{
				tea.KeyMsg{
					Type: tea.KeyEnter,
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			have := dialogcomponents.ParseTestInput("")
			want := dialogcomponents.TestInput{}
			must.Eq(t, want, have)
		})
	})

	t.Run("TestInputs.Next", func(t *testing.T) {
		t.Parallel()
		testInputs := dialogcomponents.NewTestInputs(
			dialogcomponents.TestInput{tea.KeyMsg{Type: tea.KeyCtrlA}},
			dialogcomponents.TestInput{tea.KeyMsg{Type: tea.KeyCtrlB}},
			dialogcomponents.TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}},
		)
		// request the first entry: A
		haveNext := testInputs.Next()
		wantNext := dialogcomponents.TestInput{tea.KeyMsg{Type: tea.KeyCtrlA}}
		must.Eq(t, wantNext, haveNext)
		wantRemaining := dialogcomponents.NewTestInputs(
			dialogcomponents.TestInput{tea.KeyMsg{Type: tea.KeyCtrlB}},
			dialogcomponents.TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}},
		)
		must.Eq(t, wantRemaining, testInputs)
		// request the next entry: B
		haveNext = testInputs.Next()
		wantNext = dialogcomponents.TestInput{tea.KeyMsg{Type: tea.KeyCtrlB}}
		must.Eq(t, wantNext, haveNext)
		wantRemaining = dialogcomponents.NewTestInputs(
			dialogcomponents.TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}},
		)
		must.Eq(t, wantRemaining, testInputs)
		// request the next entry: C
		haveNext = testInputs.Next()
		wantNext = dialogcomponents.TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}}
		must.Eq(t, wantNext, haveNext)
		must.EqOp(t, 0, testInputs.Len())
		// request the next entry: empty
		haveNext = testInputs.Next()
		wantNext = dialogcomponents.TestInput{}
		must.Eq(t, wantNext, haveNext)
		must.EqOp(t, 0, testInputs.Len())
	})
}
