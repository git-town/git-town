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
		keyA := TestInput{tea.KeyMsg{Type: tea.KeyCtrlA}}
		keyB := TestInput{tea.KeyMsg{Type: tea.KeyCtrlB}}
		keyC := TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}}
		testInputs := NewTestInputs(
			keyA,
			keyB,
			keyC,
		)
		must.EqOp(t, 0, testInputs.cursor)
		// request the first entry: A
		have := testInputs.Next()
		must.Eq(t, keyA, have)
		must.EqOp(t, 1, testInputs.cursor)
		must.False(t, testInputs.IsEmpty())
		// request the next entry: B
		have = testInputs.Next()
		must.Eq(t, keyB, have)
		must.False(t, testInputs.IsEmpty())
		// request the next entry: C
		have = testInputs.Next()
		must.Eq(t, keyC, have)
		must.True(t, testInputs.IsEmpty())
	})
}
