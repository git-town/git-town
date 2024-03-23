package output

import (
	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/v13/test/datatable"
)

// RenderExecutedGitCommands provides the textual Gherkin table representation of the given executed Git commands.
// The DataTable table matches the structure of the given Gherkin table.
func RenderExecutedGitCommands(commands []ExecutedGitCommand, table *messages.PickleStepArgument_PickleTable) datatable.DataTable {
	tableHasBranches := table.Rows[0].Cells[0].Value == "BRANCH"
	tableHasType := len(table.Rows[0].Cells) > 1 && table.Rows[0].Cells[1].Value == "TYPE"
	result := datatable.DataTable{}
	if tableHasBranches {
		if tableHasType {
			result.AddRow("BRANCH", "TYPE", "COMMAND")
		} else {
			result.AddRow("BRANCH", "COMMAND")
		}
	} else {
		result.AddRow("COMMAND")
	}
	lastBranch := ""
	for _, cmd := range commands {
		if tableHasBranches {
			if tableHasType {
				branch := branchForTableWithTypes(cmd, lastBranch)
				result.AddRow(branch, cmd.CommandType.String(), cmd.Command)
			} else {
				branch := branchForTableWithoutTypes(cmd, lastBranch)
				result.AddRow(branch, cmd.Command)
			}
		} else {
			result.AddRow(cmd.Command)
		}
		lastBranch = cmd.Branch
	}
	return result
}

func branchForTableWithTypes(cmd ExecutedGitCommand, lastBranch string) string {
	switch {
	case cmd.Branch == "" && cmd.CommandType == CommandTypeFrontend:
		return "<none>"
	case cmd.Branch == lastBranch:
		return ""
	case cmd.Branch == "":
		return ""
	default:
		return cmd.Branch
	}
}

func branchForTableWithoutTypes(cmd ExecutedGitCommand, lastBranch string) string {
	switch {
	case cmd.Branch == lastBranch:
		return ""
	case cmd.Branch == "":
		return "<none>"
	default:
		return cmd.Branch
	}
}
