package test

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/v8/src/stringslice"
	"github.com/git-town/git-town/v8/test/gherkin"
	"github.com/git-town/git-town/v8/test/helpers"
)

// ScenarioState constains the state that is shared by all steps within a scenario.
type ScenarioState struct {
	// the Fixture used in the current scenario
	fixture Fixture

	// initialLocalBranches contains the local branches before the WHEN steps run
	initialLocalBranches []string

	// initialRemoteBranches contains the remote branches before the WHEN steps run
	initialRemoteBranches []string

	// initialCommits describes the commits in this Git environment before the WHEN steps ran.
	initialCommits *messages.PickleStepArgument_PickleTable

	// initialBranchHierarchy describes the branch hierarchy before the WHEN steps ran.
	initialBranchHierarchy gherkin.DataTable

	// initialCurrentBranch contains the name of the branch that was checked out before the WHEN steps ran
	initialCurrentBranch string

	// the error of the last run of Git Town
	runErr error

	// indicates whether the scenario has verified the error
	runErrChecked bool

	// the output of the last run of Git Town
	runOutput string

	// content of the uncommitted file in the workspace
	uncommittedContent string

	// name of the uncommitted file in the workspace
	uncommittedFileName string
}

// Reset restores the null value of this ScenarioState.
func (state *ScenarioState) Reset(gitEnv Fixture) {
	state.fixture = gitEnv
	state.initialLocalBranches = []string{"main"}
	state.initialRemoteBranches = []string{"main"}
	state.initialCommits = nil
	state.initialBranchHierarchy = gherkin.DataTable{Cells: [][]string{{"BRANCH", "PARENT"}}}
	state.initialCurrentBranch = ""
	state.runOutput = ""
	state.runErr = nil
	state.runErrChecked = false
	state.uncommittedFileName = ""
	state.uncommittedContent = ""
}

// InitialBranches provides the branches in this Scenario before the WHEN steps ran.
func (state *ScenarioState) InitialBranches() gherkin.DataTable {
	result := gherkin.DataTable{}
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

// compareExistingCommits compares the commits in the Git environment of the given ScenarioState
// against the given Gherkin table.
func (state *ScenarioState) compareTable(table *messages.PickleStepArgument_PickleTable) error {
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
