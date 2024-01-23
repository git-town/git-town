package dialog

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/muesli/termenv"
)

const enterAliasesHelp = `
You can create shorter aliases for frequently used Git Town commands.
For example, if the "git town sync" command is aliased,
you can call it as "git sync".

Please select which Git Town commands should be shortened.
If you are not sure, select all :)

`

// Aliases lets the select which Git Town commands should have shorter aliases.
// This includes asking the user and updating the respective settings based on the user selection.
func Aliases(aliasableCommands configdomain.AliasableCommands, originalSelections configdomain.Aliases, dialogTestInput TestInput) (configdomain.Aliases, bool, error) {
	selections := NewAliasSelections(aliasableCommands, originalSelections)
	dialogData := AliasesModel{
		AllAliasableCommands: aliasableCommands,
		BubbleList:           newBubbleList(aliasableCommands.Strings(), 0),
		CurrentSelections:    selections,
		OriginalAliases:      originalSelections,
		selectedColor:        termenv.String().Foreground(termenv.ANSIGreen),
	}
	program := tea.NewProgram(dialogData)
	if len(dialogTestInput) > 0 {
		go func() {
			for _, input := range dialogTestInput {
				program.Send(input)
			}
		}()
	}
	dialogResult, err := program.Run()
	if err != nil {
		return configdomain.Aliases{}, false, err
	}
	result := dialogResult.(AliasesModel) //nolint:forcetypeassert
	if result.Status == dialogStatusAborted {
		return configdomain.Aliases{}, true, nil
	}
	selectedCommands := result.Checked(aliasableCommands)
	selectionText := DetermineAliasSelectionText(selectedCommands)
	fmt.Printf("Aliased commands: %s\n", formattedSelection(selectionText, false))
	return DetermineAliasResult(result.CurrentSelections, aliasableCommands, originalSelections), false, nil
}

type AliasesModel struct {
	BubbleList
	AllAliasableCommands configdomain.AliasableCommands
	CurrentSelections    []AliasSelection
	OriginalAliases      configdomain.Aliases
	selectedColor        termenv.Style
}

func (self AliasesModel) Checked(aliasableCommands configdomain.AliasableCommands) configdomain.AliasableCommands {
	result := configdomain.AliasableCommands{}
	for c, choice := range self.CurrentSelections {
		if choice == AliasSelectionGT {
			result = append(result, aliasableCommands[c])
		}
	}
	return result
}

func (self AliasesModel) Init() tea.Cmd {
	return nil
}

// toggleCurrentEntry unchecks the currently selected list entry if it is checked,
// and checks it if it is unchecked.
func (self *AliasesModel) RotateCurrentEntry() {
	switch self.CurrentSelections[self.Cursor] {
	case AliasSelectionNone:
		commandAtCursor := self.AllAliasableCommands[self.Cursor]
		gitTownAlias := "town " + commandAtCursor.String()
		originalAlias, hasOriginalAlias := self.OriginalAliases[commandAtCursor]
		if hasOriginalAlias && originalAlias != gitTownAlias {
			self.CurrentSelections[self.Cursor] = AliasSelectionOther
		} else {
			self.CurrentSelections[self.Cursor] = AliasSelectionGT
		}
	case AliasSelectionOther:
		self.CurrentSelections[self.Cursor] = AliasSelectionGT
	case AliasSelectionGT:
		self.CurrentSelections[self.Cursor] = AliasSelectionNone
	}
}

// SelectAll checks all entries in the list.
func (self *AliasesModel) SelectAll() {
	for s := range self.CurrentSelections {
		self.CurrentSelections[s] = AliasSelectionGT
	}
}

// SelectNone unchecks all entries in the list.
func (self *AliasesModel) SelectNone() {
	for s := range self.CurrentSelections {
		self.CurrentSelections[s] = AliasSelectionNone
	}
}

func (self AliasesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
	keyMsg, isKeyMsg := msg.(tea.KeyMsg)
	if !isKeyMsg {
		return self, nil
	}
	if handled, cmd := self.BubbleList.handleKey(keyMsg); handled {
		return self, cmd
	}
	switch keyMsg.Type { //nolint:exhaustive
	case tea.KeySpace:
		self.RotateCurrentEntry()
		return self, nil
	case tea.KeyEnter:
		self.Status = dialogStatusDone
		return self, tea.Quit
	}
	switch keyMsg.String() {
	case "a":
		self.SelectAll()
	case "n":
		self.SelectNone()
	case "o":
		self.RotateCurrentEntry()
		return self, nil
	}
	return self, nil
}

func (self AliasesModel) View() string {
	if self.Status != dialogStatusActive {
		return ""
	}
	s := strings.Builder{}
	s.WriteString(enterAliasesHelp)
	for i, branch := range self.Entries {
		selected := self.Cursor == i
		checked := self.CurrentSelections[i]
		s.WriteString(self.entryNumberStr(i))
		switch {
		case selected && checked == AliasSelectionNone:
			s.WriteString(self.Colors.selection.Styled("> [ ] " + branch))
		case selected && checked == AliasSelectionOther:
			s.WriteString(self.Colors.selection.Styled("> [o] " + branch))
		case selected && checked == AliasSelectionGT:
			s.WriteString(self.Colors.selection.Styled("> [x] " + branch))
		case !selected && checked == AliasSelectionNone:
			s.WriteString("  [ ] " + branch)
		case !selected && checked == AliasSelectionOther:
			s.WriteString("  [o] " + branch)
		case !selected && checked == AliasSelectionGT:
			s.WriteString(self.selectedColor.Styled("  [x] " + branch))
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
	// toggle
	s.WriteString(self.Colors.helpKey.Styled("space"))
	s.WriteString(self.Colors.help.Styled("/"))
	s.WriteString(self.Colors.helpKey.Styled("o"))
	s.WriteString(self.Colors.help.Styled(" toggle   "))
	// select all/none
	s.WriteString(self.Colors.helpKey.Styled("a"))
	s.WriteString(self.Colors.help.Styled("/"))
	s.WriteString(self.Colors.helpKey.Styled("n"))
	s.WriteString(self.Colors.help.Styled(" select all/none   "))
	// numbers
	s.WriteString(self.Colors.helpKey.Styled("0"))
	s.WriteString(self.Colors.help.Styled("-"))
	s.WriteString(self.Colors.helpKey.Styled("9"))
	s.WriteString(self.Colors.help.Styled(" jump   "))
	// accept
	s.WriteString(self.Colors.helpKey.Styled("enter"))
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

func DetermineAliasResult(selections []AliasSelection, allAliasableCommands configdomain.AliasableCommands, oldAliases configdomain.Aliases) configdomain.Aliases {
	result := configdomain.Aliases{}
	for s, selection := range selections {
		command := allAliasableCommands[s]
		switch selection {
		case AliasSelectionGT:
			result[command] = "town " + command.String()
		case AliasSelectionNone:
			// do nothing
		case AliasSelectionOther:
			result[command] = oldAliases[command]
		}
	}
	return result
}

func DetermineAliasSelectionText(selectedCommands configdomain.AliasableCommands) string {
	switch len(selectedCommands) {
	case 0:
		return "(none)"
	case len(configdomain.AllAliasableCommands()):
		return "(all)"
	default:
		return strings.Join(selectedCommands.Strings(), ", ")
	}
}

func NewAliasSelections(aliasableCommands configdomain.AliasableCommands, existingAliases configdomain.Aliases) []AliasSelection {
	result := make([]AliasSelection, len(aliasableCommands))
	for a, aliasableCommand := range aliasableCommands {
		existingAlias, exists := existingAliases[aliasableCommand]
		switch {
		case !exists:
			result[a] = AliasSelectionNone
		case existingAlias == "town "+aliasableCommand.String():
			result[a] = AliasSelectionGT
		default:
			result[a] = AliasSelectionOther
		}
	}
	return result
}

type AliasSelection int

const (
	AliasSelectionNone  AliasSelection = iota // no alias
	AliasSelectionGT                          // the user wants to keep the externally set alias
	AliasSelectionOther                       //
)
