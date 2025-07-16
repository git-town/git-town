package dialogcomponents_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestTestInputs(t *testing.T) {
	t.Parallel()

	t.Run("LoadTestInputs", func(t *testing.T) {
		t.Parallel()
		env := []string{
			"foo=bar",
			"GITTOWN_DIALOG_INPUT_1=enter",
			"GITTOWN_DIALOG_INPUT_2=space|down|space|5|enter",
			"GITTOWN_DIALOG_INPUT_3=ctrl+c",
		}
		have := dialogcomponents.LoadTestInputs(env)
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
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			keyA := dialogcomponents.TestInput{tea.KeyMsg{Type: tea.KeyCtrlA}}
			keyB := dialogcomponents.TestInput{tea.KeyMsg{Type: tea.KeyCtrlB}}
			keyC := dialogcomponents.TestInput{tea.KeyMsg{Type: tea.KeyCtrlC}}
			testInputs := dialogcomponents.NewTestInputs(
				keyA,
				keyB,
				keyC,
			)
			// request the first entry: A
			have := testInputs.Next()
			must.Eq(t, Some(keyA), have)
			must.False(t, testInputs.IsEmpty())
			// request the next entry: B
			have = testInputs.Next()
			must.Eq(t, Some(keyB), have)
			must.False(t, testInputs.IsEmpty())
			// request the next entry: C
			have = testInputs.Next()
			must.Eq(t, Some(keyC), have)
			must.True(t, testInputs.IsEmpty())
		})
		t.Run("not populated", func(t *testing.T) {
			t.Parallel()
			testInputs := dialogcomponents.NewTestInputs()
			// request the first entry: A
			have := testInputs.Next()
			must.Eq(t, None[dialogcomponents.TestInput](), have)
			must.True(t, testInputs.IsEmpty())
		})
	})
}
