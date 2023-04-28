package helpers

import (
	"strings"

	"github.com/cucumber/messages-go/v10"
)

func TableToInput(table *messages.PickleStepArgument_PickleTable) []string {
	var result []string
	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i]
		answer := row.Cells[1].Value
		answer = strings.ReplaceAll(answer, "[ENTER]", "\n")
		answer = strings.ReplaceAll(answer, "[DOWN]", "\x1b[B")
		answer = strings.ReplaceAll(answer, "[UP]", "\x1b[A")
		answer = strings.ReplaceAll(answer, "[SPACE]", " ")
		result = append(result, answer)
	}
	return result
}
