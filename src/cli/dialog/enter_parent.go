package dialog

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/gohacks"
	"github.com/muesli/termenv"
)

const PerennialBranchOption = "<none> (perennial branch)"

// EnterParent lets the user select the parent branch for the given branch.
func EnterParent(args EnterParentArgs) (gitdomain.LocalBranchName, bool, error) {
	parentCandidates := EnterParentEntries(args)
	numberLen := gohacks.NumberLength(len(parentCandidates))
	dialogData := enterParentModel{
		bubbleList:   newBubbleList(parentCandidates, args.MainBranch.String()),
		branch:       args.Branch.String(),
		entryNumber:  "",
		maxDigits:    numberLen,
		mainBranch:   args.MainBranch.String(),
		numberFormat: fmt.Sprintf("%%0%dd", numberLen),
	}
	dialogResult, err := tea.NewProgram(dialogData).Run()
	if err != nil {
		return gitdomain.EmptyLocalBranchName(), false, err
	}
	result := dialogResult.(enterParentModel) //nolint:forcetypeassert // we know the type for sure here
	selectedBranch := gitdomain.LocalBranchName(result.selectedEntry())
	return selectedBranch, result.aborted, nil
}

type EnterParentArgs struct {
	Branch        gitdomain.LocalBranchName
	LocalBranches gitdomain.LocalBranchNames
	Lineage       configdomain.Lineage
	MainBranch    gitdomain.LocalBranchName
}

type enterParentModel struct {
	bubbleList
	branch       string // the branch for which to enter the parent
	entryNumber  string // the currently entered branch number
	mainBranch   string // name of the main branch
	maxDigits    int    // the maximal number of digits in the branch number
	numberFormat string // template for formatting the number
}

func (self enterParentModel) Init() tea.Cmd {
	return nil
}

func (self enterParentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
	keyMsg, isKeyMsg := msg.(tea.KeyMsg)
	if !isKeyMsg {
		return self, nil
	}
	if handled, cmd := self.bubbleList.handleKey(keyMsg); handled {
		return self, cmd
	}
	if keyMsg.Type == tea.KeyEnter {
		return self, tea.Quit
	}
	switch keyStr := keyMsg.String(); keyStr {
	case "o":
		return self, tea.Quit
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		self.entryNumber += keyStr
		if len(self.entryNumber) > self.maxDigits {
			self.entryNumber = self.entryNumber[1:]
		}
		number64, _ := strconv.ParseInt(self.entryNumber, 10, 0)
		number := int(number64)
		if number < len(self.entries) {
			self.cursor = number
		}
	}
	return self, nil
}

func (self enterParentModel) View() string {
	s := strings.Builder{}
	s.WriteString("\nPlease select the parent of branch \"" + self.branch + "\" or enter its number.\n")
	s.WriteString("Most of the time this is the main development branch (" + self.mainBranch + ").\n\n")
	dim := termenv.String().Faint()
	for i, branch := range self.entries {
		if i == self.cursor {
			s.WriteString(dim.Styled(fmt.Sprintf(self.numberFormat, i)))
			s.WriteString(self.colors.selection.Styled(" > " + branch))
		} else {
			s.WriteString(dim.Styled(fmt.Sprintf(self.numberFormat, i)))
			s.WriteString("   " + branch)
		}
		s.WriteRune('\n')
	}
	s.WriteString("\n\n  ")
	// up
	s.WriteString(self.colors.helpKey.Styled("↑"))
	s.WriteString(self.colors.help.Styled("/"))
	s.WriteString(self.colors.helpKey.Styled("k"))
	s.WriteString(self.colors.help.Styled(" up   "))
	// down
	s.WriteString(self.colors.helpKey.Styled("↓"))
	s.WriteString(self.colors.help.Styled("/"))
	s.WriteString(self.colors.helpKey.Styled("j"))
	s.WriteString(self.colors.help.Styled(" down   "))
	// accept
	s.WriteString(self.colors.helpKey.Styled("enter"))
	s.WriteString(self.colors.help.Styled("/"))
	s.WriteString(self.colors.helpKey.Styled("o"))
	s.WriteString(self.colors.help.Styled(" accept   "))
	// numbers
	s.WriteString(self.colors.helpKey.Styled("0"))
	s.WriteString(self.colors.help.Styled("-"))
	s.WriteString(self.colors.helpKey.Styled("9"))
	s.WriteString(self.colors.help.Styled(" jump   "))
	// abort
	s.WriteString(self.colors.helpKey.Styled("ctrl-c"))
	s.WriteString(self.colors.help.Styled("/"))
	s.WriteString(self.colors.helpKey.Styled("q"))
	s.WriteString(self.colors.help.Styled(" abort"))
	return s.String()
}

func EnterParentEntries(args EnterParentArgs) []string {
	parentCandidateBranches := args.LocalBranches.Remove(args.Branch).Remove(args.Lineage.Children(args.Branch)...)
	parentCandidateBranches.Sort()
	parentCandidates := parentCandidateBranches.Hoist(args.MainBranch).Strings()
	return append([]string{PerennialBranchOption}, parentCandidates...)
}
