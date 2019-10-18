package steps

import (
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/test"
)

// scenarioState constains the state that is shared by all steps within a scenario.
type scenarioState struct {
	// the GitEnvironment used in the current scenario
	gitEnvironment *test.GitEnvironment

	// the result of the last run of Git Town
	lastRunOutput string

	// the error of the last run of Git Town
	lastRunErr error

	// originalCommitTable describes the commits in this Git environment before the WHEN steps ran.
	originalCommitTable *gherkin.DataTable
}
