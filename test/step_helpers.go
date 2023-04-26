package test

import (
	"fmt"

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
