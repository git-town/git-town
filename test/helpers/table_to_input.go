package helpers

import (
	"errors"
	"strings"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages/go/v21"
)

func TableToInputEnv(table *godog.Table) ([]string, error) {
	result := make([]string, 0, len(table.Rows)-1)
	keyColumn, err := detectKeysColumn(table.Rows[0])
	if err != nil {
		return result, err
	}
	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i]
		answersCucumberStyle := row.Cells[keyColumn].Value
		answersEnvStyle := strings.ReplaceAll(answersCucumberStyle, " ", "|")
		result = append(result, answersEnvStyle)
	}
	return result, nil
}

func detectKeysColumn(row *messages.PickleTableRow) (int, error) {
	for i, cell := range row.Cells {
		if cell.Value == "KEYS" {
			return i, nil
		}
	}
	return 0, errors.New(`no table column with header "KEYS" detected`)
}
