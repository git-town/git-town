package dialogcomponents_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestInputs(t *testing.T) {
	t.Parallel()

	t.Run("Inputs.Next", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			keyA := dialogcomponents.Input{
				Messages: []tea.Msg{
					tea.KeyMsg{Type: tea.KeyCtrlA},
				},
			}
			keyB := dialogcomponents.Input{
				Messages: []tea.Msg{
					tea.KeyMsg{Type: tea.KeyCtrlB},
				},
			}
			keyC := dialogcomponents.Input{
				Messages: []tea.Msg{
					tea.KeyMsg{Type: tea.KeyCtrlC},
				},
			}
			inputs := dialogcomponents.NewInputs(
				keyA,
				keyB,
				keyC,
			)
			// request the first entry: A
			have := inputs.Next()
			must.True(t, have.EqualSome(keyA))
			must.False(t, inputs.IsEmpty())
			// request the next entry: B
			have = inputs.Next()
			must.True(t, have.EqualSome(keyB))
			must.False(t, inputs.IsEmpty())
			// request the next entry: C
			have = inputs.Next()
			must.True(t, have.EqualSome(keyC))
			must.True(t, inputs.IsEmpty())
		})
		t.Run("not populated", func(t *testing.T) {
			t.Parallel()
			inputs := dialogcomponents.NewInputs()
			// request the first entry: A
			have := inputs.Next()
			must.Eq(t, None[dialogcomponents.Input](), have)
			must.True(t, inputs.IsEmpty())
		})
		t.Run("exceed given inputs", func(t *testing.T) {
			t.Parallel()
			defer func() {
				if err := recover(); err == nil {
					t.Errorf("did not panic as expected")
				}
			}()
			keyA := dialogcomponents.Input{
				Messages: []tea.Msg{
					tea.KeyMsg{Type: tea.KeyCtrlA},
				},
			}
			inputs := dialogcomponents.NewInputs(keyA)
			// request the first entry
			have := inputs.Next()
			must.True(t, have.EqualSome(keyA))
			must.True(t, inputs.IsEmpty())
			// request the next entry
			_ = inputs.Next()
		})
	})

	t.Run("LoadInputs", func(t *testing.T) {
		t.Parallel()
		env := []string{
			"foo=bar",
			"GITTOWN_DIALOG_INPUT_1=welcome@enter",
			"GITTOWN_DIALOG_INPUT_2=perennial-branches@space|down|space|5|enter",
			"GITTOWN_DIALOG_INPUT_3=perennial-regex@ctrl+c",
		}
		have := dialogcomponents.LoadInputs(env)
		want := dialogcomponents.NewInputs(
			dialogcomponents.Input{
				Messages: []tea.Msg{
					tea.KeyMsg{Type: tea.KeyEnter},
				},
				StepName: "welcome",
			},
			dialogcomponents.Input{
				Messages: []tea.Msg{
					tea.KeyMsg{Type: tea.KeySpace},
					tea.KeyMsg{Type: tea.KeyDown},
					tea.KeyMsg{Type: tea.KeySpace},
					tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}},
					tea.KeyMsg{Type: tea.KeyEnter},
				},
				StepName: "perennial-branches",
			},
			dialogcomponents.Input{
				Messages: []tea.Msg{
					tea.KeyMsg{Type: tea.KeyCtrlC},
				},
				StepName: "perennial-regex",
			},
		)
		must.Eq(t, want, have)
	})
}
