package test

import (
	"fmt"
	"strings"

	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/v8/test/helpers"
)

// compareExistingCommits compares the commits in the Git environment of the given FeatureState
// against the given Gherkin table.
func compareExistingCommits(state *ScenarioState, table *messages.PickleStepArgument_PickleTable) error {
	fields := helpers.TableFields(table)
	commitTable, err := state.fixture.CommitTable(fields)
	if err != nil {
		return fmt.Errorf("cannot determine commits in the developer repo: %w", err)
	}
	diff, errorCount := commitTable.EqualGherkin(table)
	if errorCount != 0 {
		fmt.Printf("\nERROR! Found %d differences in the existing commits\n\n", errorCount)
		fmt.Println(diff)
		return fmt.Errorf("mismatching commits found, see diff above")
	}
	return nil
}

// hasTag indicates whether the given feature has a tag with the given name.
func hasTag(scenario *messages.Pickle, name string) bool {
	for _, tag := range scenario.GetTags() {
		if tag.Name == name {
			return true
		}
	}
	return false
}

func tableToInput(table *messages.PickleStepArgument_PickleTable) []string {
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
