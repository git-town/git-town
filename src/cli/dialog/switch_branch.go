package dialog

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

type Model struct {
	Branches       []string // names of all branches
	cursor         int      // 0-based number of the selected row
	InitialBranch  string   // name of the currently checked out branch
	SelectedBranch string   // name of the currently selected branch
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			m.MoveCursorUp()
		case tea.KeyDown:
			m.MoveCursorDown()
		case tea.KeyEnter:
			m.SelectedBranch = m.Branches[m.cursor]
			return m, tea.Quit
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "k":
				m.MoveCursorUp()
			case "j":
				m.MoveCursorDown()
			case "o":
				m.SelectedBranch = m.Branches[m.cursor]
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m *Model) MoveCursorUp() {
	if m.cursor > 0 {
		m.cursor--
	} else {
		m.cursor = len(m.Branches) - 1
	}
}

func (m *Model) MoveCursorDown() {
	if m.cursor < len(m.Branches)-1 {
		m.cursor++
	} else {
		m.cursor = 0
	}
}

func (m Model) View() string {
	s := strings.Builder{}
	for _, branch := range m.Branches {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress q to quit.\n"

	return s
}

func NewBuilder(lineage configdomain.Lineage) Builder {
	return Builder{
		Entries: ModalSelectEntries{},
		Lineage: lineage,
	}
}

// queryBranch lets the user select a new branch via a visual dialog.
// Indicates via `validSelection` whether the user made a valid selection.
func SwitchBranch(roots gitdomain.LocalBranchNames, selected gitdomain.LocalBranchName, lineage configdomain.Lineage) (selection gitdomain.LocalBranchName, validSelection bool, err error) {
	builder := NewBuilder(lineage)
	err = builder.CreateEntries(roots, selected)
	if err != nil {
		return gitdomain.EmptyLocalBranchName(), false, err
	}
	choice, err := ModalSelect(builder.Entries, selected.String())
	if err != nil {
		return gitdomain.EmptyLocalBranchName(), false, err
	}
	if choice == nil {
		return gitdomain.EmptyLocalBranchName(), false, nil
	}
	return gitdomain.NewLocalBranchName(*choice), true, nil
}

// Builder builds up the switch-branch dialog entries.
type Builder struct {
}

// AddEntryAndChildren adds the given branch and all its child branches to the given entries collection.
func (self *Builder) AddEntryAndChildren(branch gitdomain.LocalBranchName, indent int) error {
	self.Entries = append(self.Entries, ModalSelectEntry{
		Text:  strings.Repeat("  ", indent) + branch.String(),
		Value: branch.String(),
	})
	var err error
	for _, child := range self.Lineage.Children(branch) {
		err = self.AddEntryAndChildren(child, indent+1)
		if err != nil {
			return err
		}
	}
	return nil
}

// createEntries provides all the entries for the branch dialog.
func (self *Builder) CreateEntries(roots gitdomain.LocalBranchNames, selected gitdomain.LocalBranchName) error {
	var err error
	for _, root := range roots {
		err = self.AddEntryAndChildren(root, 0)
		if err != nil {
			return err
		}
	}
	if len(self.Entries) == 0 {
		self.Entries = append(self.Entries, ModalSelectEntry{
			Text:  string(selected),
			Value: string(selected),
		})
	}
	return nil
}
