package test

import (
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/test/gherkintools"
)

// RenderExecutedGitCommands provides the textual Gherkin table representation of the given executed Git commands.
// The Mortadella table matches the structure of the given Gherkin table.
func RenderExecutedGitCommands(commands []ExecutedGitCommand, table *gherkin.DataTable) gherkintools.Mortadella {
	tableHasBranches := table.Rows[0].Cells[0].Value == "BRANCH"
	morta := gherkintools.Mortadella{}
	if tableHasBranches {
		morta.AddRow("BRANCH", "COMMAND")
	} else {
		morta.AddRow("COMMAND")
	}
	lastBranch := ""
	for _, cmd := range commands {
		if tableHasBranches {
			if cmd.Branch == lastBranch {
				morta.AddRow("", cmd.Command)
			} else {
				morta.AddRow(cmd.Branch, cmd.Command)
			}
		} else {
			morta.AddRow(cmd.Command)
		}
		lastBranch = cmd.Branch
	}
	return morta
}

// RenderTable provides the textual Gherkin representation of the given Gherkin table.
func RenderTable(table *gherkin.DataTable) string {
	morta := gherkintools.Mortadella{}
	for _, row := range table.Rows {
		values := []string{}
		for _, cell := range row.Cells {
			values = append(values, cell.Value)
		}
		morta.AddRow(values...)
	}
	return morta.String()
}
