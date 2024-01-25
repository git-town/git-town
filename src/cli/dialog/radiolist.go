package dialog

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// EnterMainBranch lets the user select a new main branch for this repo.
func radioList[C fmt.Stringer](entries []C, cursor int, help string, testInput TestInput) (selected C, aborted bool, err error) { //nolint:ireturn
	program := tea.NewProgram(radioListModel[C]{
		BubbleList: newBubbleList(entries, cursor),
		help:       help,
	})
	if len(testInput) > 0 {
		go func() {
			for _, input := range testInput {
				program.Send(input)
			}
		}()
	}
	dialogResult, err := program.Run()
	if err != nil {
		return entries[0], false, err
	}
	result := dialogResult.(radioListModel[C]) //nolint:forcetypeassert
	return result.selectedEntry(), result.aborted(), nil
}

type radioListModel[C fmt.Stringer] struct {
	BubbleList[C]
	help string // help text to display before the radio list
}

func (self radioListModel[C]) Init() tea.Cmd {
	return nil
}

func (self radioListModel[C]) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
	keyMsg, isKeyMsg := msg.(tea.KeyMsg)
	if !isKeyMsg {
		return self, nil
	}
	if handled, cmd := self.BubbleList.handleKey(keyMsg); handled {
		return self, cmd
	}
	if keyMsg.Type == tea.KeyEnter {
		self.Status = dialogStatusDone
		return self, tea.Quit
	}
	if keyMsg.String() == "o" {
		self.Status = dialogStatusDone
		return self, tea.Quit
	}
	return self, nil
}

func (self radioListModel[C]) View() string {
	if self.Status != dialogStatusActive {
		return ""
	}
	s := strings.Builder{}
	s.WriteString(self.help)
	for i, branch := range self.Entries {
		s.WriteString(self.entryNumberStr(i))
		if i == self.Cursor {
			s.WriteString(self.Colors.selection.Styled("> " + branch.String()))
		} else {
			s.WriteString("  " + branch.String())
		}
		s.WriteRune('\n')
	}
	s.WriteString("\n\n  ")
	// up
	s.WriteString(self.Colors.helpKey.Styled("↑"))
	s.WriteString(self.Colors.help.Styled("/"))
	s.WriteString(self.Colors.helpKey.Styled("k"))
	s.WriteString(self.Colors.help.Styled(" up   "))
	// down
	s.WriteString(self.Colors.helpKey.Styled("↓"))
	s.WriteString(self.Colors.help.Styled("/"))
	s.WriteString(self.Colors.helpKey.Styled("j"))
	s.WriteString(self.Colors.help.Styled(" down   "))
	// numbers
	s.WriteString(self.Colors.helpKey.Styled("0"))
	s.WriteString(self.Colors.help.Styled("-"))
	s.WriteString(self.Colors.helpKey.Styled("9"))
	s.WriteString(self.Colors.help.Styled(" jump   "))
	// accept
	s.WriteString(self.Colors.helpKey.Styled("enter"))
	s.WriteString(self.Colors.help.Styled("/"))
	s.WriteString(self.Colors.helpKey.Styled("o"))
	s.WriteString(self.Colors.help.Styled(" accept   "))
	// abort
	s.WriteString(self.Colors.helpKey.Styled("q"))
	s.WriteString(self.Colors.help.Styled("/"))
	s.WriteString(self.Colors.helpKey.Styled("esc"))
	s.WriteString(self.Colors.help.Styled("/"))
	s.WriteString(self.Colors.helpKey.Styled("ctrl-c"))
	s.WriteString(self.Colors.help.Styled(" abort"))
	return s.String()
}
