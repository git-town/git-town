package test

import (
	"sort"
	"strings"

	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/v7/src/run"
	"github.com/git-town/git-town/v7/src/stringslice"
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

	// initialCurrentBranch contains the name of the branch that was checked out before the WHEN steps ran
	initialCurrentBranch string

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
	state.initialRemoteBranches = []string{"main"}
	state.initialCommits = nil
	state.initialBranchHierarchy = DataTable{Cells: [][]string{{"BRANCH", "PARENT"}}}
	state.initialCurrentBranch = ""
	state.runRes = nil
	state.runErr = nil
	state.runErrChecked = false
	state.uncommittedFileName = ""
	state.uncommittedContent = ""
}

// InitialBranches provides the branches in this Scenario before the WHEN steps ran.
func (state *ScenarioState) InitialBranches() DataTable {
	result := DataTable{}
	result.AddRow("REPOSITORY", "BRANCHES")
	sort.Strings(state.initialLocalBranches)
	state.initialLocalBranches = stringslice.Hoist(state.initialLocalBranches, "main")
	sort.Strings(state.initialRemoteBranches)
	state.initialRemoteBranches = stringslice.Hoist(state.initialRemoteBranches, "main")
	localBranchesJoined := strings.Join(state.initialLocalBranches, ", ")
	remoteBranchesJoined := strings.Join(state.initialRemoteBranches, ", ")
	if localBranchesJoined == remoteBranchesJoined {
		result.AddRow("local, origin", localBranchesJoined)
	} else {
		result.AddRow("local", localBranchesJoined)
		if remoteBranchesJoined != "" {
			result.AddRow("origin", remoteBranchesJoined)
		}
	}
	return result
}
