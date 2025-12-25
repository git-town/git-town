package dialogcomponents_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/shoenig/test/must"
)

func TestInput(t *testing.T) {
	t.Parallel()

	t.Run("ParseInput", func(t *testing.T) {
		t.Parallel()
		t.Run("multiple values", func(t *testing.T) {
			t.Parallel()
			have := dialogcomponents.ParseInput("step@enter|space|ctrl+c")
			want := dialogcomponents.Input{
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
				StepName: "step",
			}
			must.Eq(t, want, have)
		})
		t.Run("single value", func(t *testing.T) {
			t.Parallel()
			have := dialogcomponents.ParseInput("step@enter")
			want := dialogcomponents.Input{
				Messages: []tea.Msg{
					tea.KeyMsg{
						Type: tea.KeyEnter,
					},
				},
				StepName: "step",
			}
			must.Eq(t, want, have)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			have := dialogcomponents.ParseInput("step@")
			want := dialogcomponents.Input{
				Messages: []tea.Msg{},
				StepName: "step",
			}
			must.Eq(t, want, have)
		})
	})
}
