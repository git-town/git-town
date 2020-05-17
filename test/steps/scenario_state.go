package steps

import (
	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/src/command"
	"github.com/git-town/git-town/test"
)

// scenarioState constains the state that is shared by all steps within a scenario.
type scenarioState struct {
	// the GitEnvironment used in the current scenario
	gitEnv *test.GitEnvironment

	// the error of the last run of Git Town
	runErr error

	// the outcome of the last run of Git Town
	runRes *command.Result

	// initialCommits describes the commits in this Git environment before the WHEN steps ran.
	initialCommits *messages.PickleStepArgument_PickleTable

	// name of the uncommitted file in the workspace
	uncommittedFileName string

	// content of the uncommitted file in the workspace
	uncommittedContent string
}
