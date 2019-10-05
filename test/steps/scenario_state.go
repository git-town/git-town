package steps

import "github.com/Originate/git-town/test"

// scenarioState constains the state that is shared by all steps within a scenario.
type scenarioState struct {
	// the GitEnvironment used in the current scenario
	gitEnvironment *test.GitEnvironment

	// the result of the last run of Git Town
	lastRunOutput string

	// the error of the last run of Git Town
	lastRunErr error
}
