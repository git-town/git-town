package dialogcomponents_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/shoenig/test/must"
)

func TestTestInput(t *testing.T) {
	t.Parallel()

	t.Run("ParseTestInput", func(t *testing.T) {
		t.Parallel()
		t.Run("multiple values", func(t *testing.T) {
			t.Parallel()
			have := dialogcomponents.ParseTestInput("enter|space|ctrl+c")
			want := dialogcomponents.TestInput{
				Messages: []tea.Msg{
					tea.KeyMsg{
						Type: tea.KeyEnter,
					},
					tea.KeyMsg{
						Type: tea.KeySpace,
					},
					tea.KeyMsg{
						Type: tea.KeyCtrlC,
					},
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("single value", func(t *testing.T) {
			t.Parallel()
			have := dialogcomponents.ParseTestInput("enter")
			want := dialogcomponents.TestInput{
				Messages: []tea.Msg{
					tea.KeyMsg{
						Type: tea.KeyEnter,
					},
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			have := dialogcomponents.ParseTestInput("")
			want := dialogcomponents.TestInput{
				Messages: []tea.Msg{},
			}
			must.Eq(t, want, have)
		})
	})
}
