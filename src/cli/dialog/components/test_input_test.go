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
		want := components.TestInputs{
			components.TestInput{tea.KeyMsg{Type: tea.KeyEnter}}, //nolint:exhaustruct
			components.TestInput{
				tea.KeyMsg{Type: tea.KeySpace},                     //nolint:exhaustruct
				tea.KeyMsg{Type: tea.KeyDown},                      //nolint:exhaustruct
				tea.KeyMsg{Type: tea.KeySpace},                     //nolint:exhaustruct
				tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}}, //nolint:exhaustruct
				tea.KeyMsg{Type: tea.KeyEnter},                     //nolint:exhaustruct
			},
			components.TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}}, //nolint:exhaustruct
		}
		must.Eq(t, want, have)
	})

	t.Run("TestInputs.Next", func(t *testing.T) {
		t.Parallel()
		testInputs := components.TestInputs{
			components.TestInput{tea.KeyMsg{Type: tea.KeyCtrlA}}, //nolint:exhaustruct
			components.TestInput{tea.KeyMsg{Type: tea.KeyCtrlB}}, //nolint:exhaustruct
			components.TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}}, //nolint:exhaustruct
		}
		// request the first entry: A
		haveNext := testInputs.Next()
		wantNext := components.TestInput{tea.KeyMsg{Type: tea.KeyCtrlA}} //nolint:exhaustruct
		must.Eq(t, wantNext, haveNext)
		wantRemaining := components.TestInputs{
			components.TestInput{tea.KeyMsg{Type: tea.KeyCtrlB}}, //nolint:exhaustruct
			components.TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}}, //nolint:exhaustruct
		}
		must.Eq(t, wantRemaining, testInputs)
		// request the next entry: B
		haveNext = testInputs.Next()
		wantNext = components.TestInput{tea.KeyMsg{Type: tea.KeyCtrlB}} //nolint:exhaustruct
		must.Eq(t, wantNext, haveNext)
		wantRemaining = components.TestInputs{
			components.TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}}, //nolint:exhaustruct
		}
		must.Eq(t, wantRemaining, testInputs)
		// request the next entry: C
		haveNext = testInputs.Next()
		wantNext = components.TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}} //nolint:exhaustruct
		must.Eq(t, wantNext, haveNext)
		wantRemaining = components.TestInputs{}
		must.Eq(t, wantRemaining, testInputs)
		// request the next entry: empty
		haveNext = testInputs.Next()
		wantNext = components.TestInput{}
		must.Eq(t, wantNext, haveNext)
		wantRemaining = components.TestInputs{}
		must.Eq(t, wantRemaining, testInputs)
	})

	t.Run("ParseTestInput", func(t *testing.T) {
		t.Parallel()
		t.Run("multiple values", func(t *testing.T) {
			t.Parallel()
			have := components.ParseTestInput("enter|space|ctrl+c")
			want := components.TestInput{
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
			have := components.ParseTestInput("enter")
			want := components.TestInput{
				tea.KeyMsg{ //nolint:exhaustruct
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
