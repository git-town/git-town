package steps

import (
	"fmt"

	"github.com/Originate/git-town/test/helpers"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/gherkin"
)

// CommitSteps defines Cucumber step implementations around configuration.
func CommitSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^my repository is left with my original commits$`, func() error {
		return compareExistingCommits(fs, fs.activeScenarioState.originalCommitTable)
	})

	suite.Step(`^my repository now has the following commits$`, func(table *gherkin.DataTable) error {
		return compareExistingCommits(fs, table)
	})

	suite.Step(`^the following commits exist in my repository$`, func(table *gherkin.DataTable) error {
		fs.activeScenarioState.originalCommitTable = table
		return fs.activeScenarioState.gitEnvironment.CreateCommits(table)
	})
}

// compareExistingCommits compares the commits in the Git environment of the given FeatureState
// against the given Gherkin table.
func compareExistingCommits(fs *FeatureState, table *gherkin.DataTable) error {
	fields := helpers.TableFields(table)
	commitTable, err := fs.activeScenarioState.gitEnvironment.CommitTable(fields)
	if err != nil {
		return fmt.Errorf("cannot determine commits in the developer repo: %w", err)
	}
	diff, errorCount := commitTable.Equal(table)
	if errorCount != 0 {
		return fmt.Errorf("found %d differences:\n%s", errorCount, diff)
	}
	return nil
}
