package helpers

import (
	"strings"

	"github.com/cucumber/godog"
	messages "github.com/cucumber/messages/go/v21"
)

func TableToInputEnv(table *godog.Table) []string {
	result := make([]string, 0, len(table.Rows)-1)
	keyColumn := detectKeysColumn(table.Rows[0])
	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i]
		answersEnvStyle := strings.ReplaceAll(row.Cells[keyColumn].Value, " ", "|")
		if len(answersEnvStyle) > 0 {
			result = append(result, answersEnvStyle)
		}
	}
	return result
}

func detectKeysColumn(row *messages.PickleTableRow) int {
	for i, cell := range row.Cells {
		if cell.Value == "KEYS" {
			return i
		}
	}
	panic(`no table column with header "KEYS" detected`)
}
