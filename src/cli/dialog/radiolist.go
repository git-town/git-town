package dialog

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// EnterMainBranch lets the user select a new main branch for this repo.
func radioList(args radioListArgs) (selected string, aborted bool, err error) {
	model := radioListModel{
		BubbleList: newBubbleList(args.entries, DetermineCursorPos(args.entries, args.defaultEntry)),
		help:       args.help,
	}
	program := tea.NewProgram(model)
	if len(args.testInput) > 0 {
		go func() {
			for _, input := range args.testInput {
				program.Send(input)
			}
		}()
	}
	dialogResult, err := program.Run()
	if err != nil {
		return "", false, err
	}
	result := dialogResult.(radioListModel) //nolint:forcetypeassert
	return result.selectedEntry(), result.Status == dialogStatusAborted, nil
}

type radioListArgs struct {
	entries      []string
	defaultEntry string
	help         string
	testInput    TestInput
}

type radioListModel struct {
	BubbleList
	help string // help text to display before the radio list
}

func (self radioListModel) Init() tea.Cmd {
	return nil
}

func (self radioListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
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

func (self radioListModel) View() string {
	if self.Status != dialogStatusActive {
		return ""
	}
	s := strings.Builder{}
	s.WriteString(self.help)
	for i, branch := range self.Entries {
		s.WriteString(self.entryNumberStr(i))
		if i == self.Cursor {
			s.WriteString(self.Colors.selection.Styled("> " + branch))
		} else {
			s.WriteString("  " + branch)
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
