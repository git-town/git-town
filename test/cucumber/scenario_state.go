package cucumber

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/v9/src/slice"
	"github.com/git-town/git-town/v9/test/datatable"
	"github.com/git-town/git-town/v9/test/fixture"
	"github.com/git-town/git-town/v9/test/helpers"
)

// ScenarioState constains the state that is shared by all steps within a scenario.
type ScenarioState struct {
	// the Fixture used in the current scenario
	fixture fixture.Fixture

	// initialLocalBranches contains the local branches before the WHEN steps run
	initialLocalBranches []string

	// initialRemoteBranches contains the remote branches before the WHEN steps run
	initialRemoteBranches []string

	// initialCommits describes the commits in this Git environment before the WHEN steps ran.
	initialCommits *messages.PickleStepArgument_PickleTable

	// initialBranchHierarchy describes the branch hierarchy before the WHEN steps ran.
	initialBranchHierarchy datatable.DataTable

	// initialCurrentBranch contains the name of the branch that was checked out before the WHEN steps ran
	initialCurrentBranch string

	// the error of the last run of Git Town
	runExitCode int

	// indicates whether the scenario has verified the error
	runExitCodeChecked bool

	// the output of the last run of Git Town
	runOutput string

	// content of the uncommitted file in the workspace
	uncommittedContent string

	// name of the uncommitted file in the workspace
	uncommittedFileName string
}

// Reset restores the null value of this ScenarioState.
func (state *ScenarioState) Reset(gitEnv fixture.Fixture) {
	state.fixture = gitEnv
	state.initialLocalBranches = []string{"main"}
	state.initialRemoteBranches = []string{"main"}
	state.initialCommits = nil
	state.initialBranchHierarchy = datatable.DataTable{Cells: [][]string{{"BRANCH", "PARENT"}}}
	state.initialCurrentBranch = ""
	state.runOutput = ""
	state.runExitCode = 0
	state.runExitCodeChecked = false
	state.uncommittedFileName = ""
	state.uncommittedContent = ""
}

// InitialBranches provides the branches in this Scenario before the WHEN steps ran.
func (state *ScenarioState) InitialBranches() datatable.DataTable {
	result := datatable.DataTable{}
	result.AddRow("REPOSITORY", "BRANCHES")
	sort.Strings(state.initialLocalBranches)
	state.initialLocalBranches = slice.Hoist(state.initialLocalBranches, "main")
	sort.Strings(state.initialRemoteBranches)
	state.initialRemoteBranches = slice.Hoist(state.initialRemoteBranches, "main")
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
	commitTable := state.fixture.CommitTable(fields)
	diff, errorCount := commitTable.EqualGherkin(table)
	if errorCount != 0 {
		fmt.Printf("\nERROR! Found %d differences in the existing commits\n\n", errorCount)
		fmt.Println(diff)
		return fmt.Errorf("mismatching commits found, see diff above")
	}
	return nil
}
