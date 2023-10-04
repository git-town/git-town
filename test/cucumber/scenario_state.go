package cucumber

import (
	"fmt"

	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/v9/src/domain"
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
	initialLocalBranches domain.LocalBranchNames

	// initialRemoteBranches contains the remote branches before the WHEN steps run
	initialRemoteBranches domain.LocalBranchNames // the remote branches are tracked as local branches in the remote repo

	// initialDevSHAs is only for looking up SHAs that existed at the developer repo before the first Git Town command ran.
	// It's not a source of truth for which branches existed at that time
	// because it might contain non-existing remote branches or miss existing remote branches.
	// An example is when origin removes a branch. initialDevSHAs will still list it
	// because the developer workspace hasn't fetched updates yet.
	initialDevSHAs map[string]domain.SHA

	// initialOriginSHAs is only for looking up SHAs that existed at the origin repo before the first Git Town command was run.
	initialOriginSHAs map[string]domain.SHA

	// initialCommits describes the commits in this Git environment before the WHEN steps ran.
	initialCommits *messages.PickleStepArgument_PickleTable

	// initialBranchHierarchy describes the branch hierarchy before the WHEN steps ran.
	initialBranchHierarchy datatable.DataTable

	// initialCurrentBranch contains the name of the branch that was checked out before the WHEN steps ran
	initialCurrentBranch domain.LocalBranchName

	// insideGitRepo indicates whether the developer workspace contains a Git repository
	insideGitRepo bool

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
func (ss *ScenarioState) Reset(gitEnv fixture.Fixture) {
	ss.fixture = gitEnv
	ss.initialLocalBranches = domain.NewLocalBranchNames("main")
	ss.initialRemoteBranches = domain.NewLocalBranchNames("main")
	ss.initialDevSHAs = map[string]domain.SHA{}
	ss.initialOriginSHAs = map[string]domain.SHA{}
	ss.initialBranchHierarchy = datatable.DataTable{Cells: [][]string{{"BRANCH", "PARENT"}}}
	ss.initialCurrentBranch = domain.EmptyLocalBranchName()
	ss.insideGitRepo = true
	ss.runOutput = ""
	ss.runExitCode = 0
	ss.runExitCodeChecked = false
	ss.uncommittedFileName = ""
	ss.uncommittedContent = ""
}

// InitialBranches provides the branches in this Scenario before the WHEN steps ran.
func (ss *ScenarioState) InitialBranches() datatable.DataTable {
	result := datatable.DataTable{}
	result.AddRow("REPOSITORY", "BRANCHES")
	ss.initialLocalBranches.Sort()
	ss.initialLocalBranches = slice.Hoist(ss.initialLocalBranches, domain.NewLocalBranchName("main"))
	ss.initialRemoteBranches.Sort()
	ss.initialRemoteBranches = slice.Hoist(ss.initialRemoteBranches, domain.NewLocalBranchName("main"))
	localBranchesJoined := ss.initialLocalBranches.Join(", ")
	remoteBranchesJoined := ss.initialRemoteBranches.Join(", ")
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
func (ss *ScenarioState) compareTable(table *messages.PickleStepArgument_PickleTable) error {
	fields := helpers.TableFields(table)
	commitTable := ss.fixture.CommitTable(fields)
	diff, errorCount := commitTable.EqualGherkin(table)
	if errorCount != 0 {
		fmt.Printf("\nERROR! Found %d differences in the existing commits\n\n", errorCount)
		fmt.Println(diff)
		return fmt.Errorf("mismatching commits found, see diff above")
	}
	return nil
}
