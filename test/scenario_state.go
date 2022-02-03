package test

import (
	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/v7/src/run"
)

// ScenarioState constains the state that is shared by all steps within a scenario.
type ScenarioState struct {
	// the GitEnvironment used in the current scenario
	gitEnv GitEnvironment

	// initialLocalBranches contains the local branches before the WHEN steps run
	initialLocalBranches []string

	// initialRemoteBranches contains the remote branches before the WHEN steps run
	initialRemoteBranches []string

	// initialCommits describes the commits in this Git environment before the WHEN steps ran.
	initialCommits *messages.PickleStepArgument_PickleTable

	// initialBranchHierarchy describes the branch hierarchy before the WHEN steps ran.
	initialBranchHierarchy DataTable

	// the error of the last run of Git Town
	runErr error

	// indicates whether the scenario has verified the error
	runErrChecked bool

	// the outcome of the last run of Git Town
	runRes *run.Result

	// content of the uncommitted file in the workspace
	uncommittedContent string

	// name of the uncommitted file in the workspace
	uncommittedFileName string
}

// Reset restores the null value of this ScenarioState.
func (state *ScenarioState) Reset(gitEnv GitEnvironment) {
	state.gitEnv = gitEnv
	state.initialLocalBranches = []string{"main"}
	state.initialRemoteBranches = []string{}
	state.initialCommits = nil
	state.initialBranchHierarchy = DataTable{Cells: [][]string{{"BRANCH", "PARENT"}}}
	state.runRes = nil
	state.runErr = nil
	state.runErrChecked = false
	state.uncommittedFileName = ""
	state.uncommittedContent = ""
}
