package test

import (
	"github.com/DATA-DOG/godog/gherkin"
)

// RenderExecutedGitCommands provides the textual Gherkin table representation of the given executed Git commands.
// The DataTable table matches the structure of the given Gherkin table.
func RenderExecutedGitCommands(commands []ExecutedGitCommand, table *gherkin.DataTable) DataTable {
	tableHasBranches := table.Rows[0].Cells[0].Value == "BRANCH"
	morta := DataTable{}
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
	morta := DataTable{}
	for _, row := range table.Rows {
		values := []string{}
		for _, cell := range row.Cells {
			values = append(values, cell.Value)
		}
		morta.AddRow(values...)
	}
	return morta.String()
}
