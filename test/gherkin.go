package test

import (
	"fmt"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/dchest/uniuri"
)

// CommitTableEntry contains the elements of a Gherkin table defining commit data.
type CommitTableEntry struct {
	branch      string
	location    string
	message     string
	fileName    string
	fileContent string
}

// NewCommitTableEntry provides a new CommitTableEntry with default values
func NewCommitTableEntry() CommitTableEntry {
	return CommitTableEntry{
		fileName:    "default_file_name_" + uniuri.NewLen(10),
		message:     "default commit message",
		location:    "local and remote",
		branch:      "main",
		fileContent: "default file content",
	}
}

// RenderSlice returns the Gherkin table representation of the given slice
func RenderSlice(data []string) string {

	// determine the width of the longest string
	width := 0
	for _, text := range data {
		if len(text) > width {
			width = len(text)
		}
	}

	// render
	result := ""
	formatStr := fmt.Sprintf("| %%-%dv |\n", width)
	for _, text := range data {
		result += fmt.Sprintf(formatStr, text)
	}
	return result
}

// RenderTable returns the Gherkin representation of the given Gherkin table
func RenderTable(table *gherkin.DataTable) string {

	// determine the width of the table
	width := 0
	for _, row := range table.Rows {
		cellWidth := len(row.Cells[0].Value)
		if (cellWidth) > width {
			width = cellWidth
		}
	}

	// convert table to slice
	slice := []string{}
	for i := 1; i < len(table.Rows); i++ {
		slice = append(slice, table.Rows[i].Cells[0].Value)
	}

	return RenderSlice(slice)
}
