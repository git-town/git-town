package helpers

import (
	"strings"

	"github.com/cucumber/messages-go/v10"
)

func TableToInputEnv(table *messages.PickleStepArgument_PickleTable) string {
	answers := make([]string, 0, len(table.Rows))
	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i]
		answer := row.Cells[0].Value
		answers = append(answers, answer)
	}
	return strings.Join(answers, "|")
}
