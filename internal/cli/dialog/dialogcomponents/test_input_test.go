package dialogcomponents

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
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
		have := LoadTestInputs(give)
		want := NewTestInputs(
			TestInput{tea.KeyMsg{Type: tea.KeyEnter}},
			TestInput{
				tea.KeyMsg{Type: tea.KeySpace},
				tea.KeyMsg{Type: tea.KeyDown},
				tea.KeyMsg{Type: tea.KeySpace},
				tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}},
				tea.KeyMsg{Type: tea.KeyEnter},
			},
			TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}},
		)
		must.Eq(t, want, have)
	})

	t.Run("ParseTestInput", func(t *testing.T) {
		t.Parallel()
		t.Run("multiple values", func(t *testing.T) {
			t.Parallel()
			have := parseTestInput("enter|space|ctrl+c")
			want := TestInput{
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
			have := parseTestInput("enter")
			want := TestInput{
				tea.KeyMsg{
					Type: tea.KeyEnter,
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			have := parseTestInput("")
			want := TestInput{}
			must.Eq(t, want, have)
		})
	})

	t.Run("TestInputs.Next", func(t *testing.T) {
		t.Parallel()
		testInputs := NewTestInputs(
			TestInput{tea.KeyMsg{Type: tea.KeyCtrlA}},
			TestInput{tea.KeyMsg{Type: tea.KeyCtrlB}},
			TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}},
		)
		// request the first entry: A
		haveNext := testInputs.Next()
		wantNext := TestInput{tea.KeyMsg{Type: tea.KeyCtrlA}}
		must.Eq(t, wantNext, haveNext)
		must.False(t, testInputs.IsEmpty())
		// request the next entry: B
		haveNext = testInputs.Next()
		wantNext = TestInput{tea.KeyMsg{Type: tea.KeyCtrlB}}
		must.Eq(t, wantNext, haveNext)
		must.False(t, testInputs.IsEmpty())
		// request the next entry: C
		haveNext = testInputs.Next()
		wantNext = TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}}
		must.Eq(t, wantNext, haveNext)
		must.True(t, testInputs.IsEmpty())
	})
}
