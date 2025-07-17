package helpers

import (
	"fmt"
	"strings"

	"github.com/cucumber/godog"
	messages "github.com/cucumber/messages/go/v21"
)

func TableToInputEnv(table *godog.Table) []string {
	result := make([]string, 0, len(table.Rows)-1)
	dialogColumn := detectColumn("DIALOG", table.Rows[0])
	inputColumn := detectColumn("KEYS", table.Rows[0])
	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i]
		dialogName := strings.ReplaceAll(row.Cells[dialogColumn].Value, " ", "-")
		input := strings.ReplaceAll(row.Cells[inputColumn].Value, " ", "|")
		if len(input) > 0 {
			result = append(result, dialogName+"@"+input)
		}
	}
	return result
}

func detectColumn(caption string, row *messages.PickleTableRow) int {
	for i, cell := range row.Cells {
		if cell.Value == caption {
			return i
		}
	}
	panic(fmt.Sprintf("no table column with header %q detected", caption))
}
