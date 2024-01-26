package dialogcomponents_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v11/src/cli/dialog/dialogcomponents"
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
		want := dialogcomponents.TestInputs{
			dialogcomponents.TestInput{tea.KeyMsg{Type: tea.KeyEnter}}, //nolint:exhaustruct
			dialogcomponents.TestInput{
				tea.KeyMsg{Type: tea.KeySpace},                     //nolint:exhaustruct
				tea.KeyMsg{Type: tea.KeyDown},                      //nolint:exhaustruct
				tea.KeyMsg{Type: tea.KeySpace},                     //nolint:exhaustruct
				tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}}, //nolint:exhaustruct
				tea.KeyMsg{Type: tea.KeyEnter},                     //nolint:exhaustruct
			},
			dialogcomponents.TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}}, //nolint:exhaustruct
		}
		must.Eq(t, want, have)
	})

	t.Run("TestInputs.Next", func(t *testing.T) {
		t.Parallel()
		testInputs := dialogcomponents.TestInputs{
			dialogcomponents.TestInput{tea.KeyMsg{Type: tea.KeyCtrlA}}, //nolint:exhaustruct
			dialogcomponents.TestInput{tea.KeyMsg{Type: tea.KeyCtrlB}}, //nolint:exhaustruct
			dialogcomponents.TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}}, //nolint:exhaustruct
		}
		// request the first entry: A
		haveNext := testInputs.Next()
		wantNext := dialogcomponents.TestInput{tea.KeyMsg{Type: tea.KeyCtrlA}} //nolint:exhaustruct
		must.Eq(t, wantNext, haveNext)
		wantRemaining := dialogcomponents.TestInputs{
			dialogcomponents.TestInput{tea.KeyMsg{Type: tea.KeyCtrlB}}, //nolint:exhaustruct
			dialogcomponents.TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}}, //nolint:exhaustruct
		}
		must.Eq(t, wantRemaining, testInputs)
		// request the next entry: B
		haveNext = testInputs.Next()
		wantNext = dialogcomponents.TestInput{tea.KeyMsg{Type: tea.KeyCtrlB}} //nolint:exhaustruct
		must.Eq(t, wantNext, haveNext)
		wantRemaining = dialogcomponents.TestInputs{
			dialogcomponents.TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}}, //nolint:exhaustruct
		}
		must.Eq(t, wantRemaining, testInputs)
		// request the next entry: C
		haveNext = testInputs.Next()
		wantNext = dialogcomponents.TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}} //nolint:exhaustruct
		must.Eq(t, wantNext, haveNext)
		wantRemaining = dialogcomponents.TestInputs{}
		must.Eq(t, wantRemaining, testInputs)
		// request the next entry: empty
		haveNext = testInputs.Next()
		wantNext = dialogcomponents.TestInput{}
		must.Eq(t, wantNext, haveNext)
		wantRemaining = dialogcomponents.TestInputs{}
		must.Eq(t, wantRemaining, testInputs)
	})

	t.Run("ParseTestInput", func(t *testing.T) {
		t.Parallel()
		t.Run("multiple values", func(t *testing.T) {
			t.Parallel()
			have := dialogcomponents.ParseTestInput("enter|space|ctrl+c")
			want := dialogcomponents.TestInput{
				tea.KeyMsg{ //nolint:exhaustruct
					Type: tea.KeyEnter,
				},
				tea.KeyMsg{ //nolint:exhaustruct
					Type: tea.KeySpace,
				},
				tea.KeyMsg{ //nolint:exhaustruct
					Type: tea.KeyCtrlC,
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("single value", func(t *testing.T) {
			t.Parallel()
			have := dialogcomponents.ParseTestInput("enter")
			want := dialogcomponents.TestInput{
				tea.KeyMsg{ //nolint:exhaustruct
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
}
