package cucumber

import (
	"errors"
	"fmt"

	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/git-town/git-town/v13/test/datatable"
	"github.com/git-town/git-town/v13/test/fixture"
	"github.com/git-town/git-town/v13/test/helpers"
)

// ScenarioState constains the state that is shared by all steps within a scenario.
type ScenarioState struct {
	// the Fixture used in the current scenario
	fixture fixture.Fixture

	// initialCommits describes the commits in this Git environment before the WHEN steps ran.
	initialCommits *datatable.DataTable

	// initialCurrentBranch contains the name of the branch that was checked out before the WHEN steps ran
	initialCurrentBranch gitdomain.LocalBranchName

	// initialDevSHAs is only for looking up SHAs that existed at the developer repo before the first Git Town command ran.
	// It's not a source of truth for which branches existed at that time
	// because it might contain non-existing remote branches or miss existing remote branches.
	// An example is when origin removes a branch. initialDevSHAs will still list it
	// because the developer workspace hasn't fetched updates yet.
	initialDevSHAs map[string]gitdomain.SHA

	// initialLineage describes the lineage before the WHEN steps ran.
	initialLineage datatable.DataTable

	// initialLocalBranches contains the local branches before the WHEN steps run
	initialLocalBranches gitdomain.LocalBranchNames

	// initialOriginSHAs is only for looking up SHAs that existed at the origin repo before the first Git Town command was run.
	initialOriginSHAs map[string]gitdomain.SHA

	// initialRemoteBranches contains the remote branches before the WHEN steps run
	initialRemoteBranches gitdomain.LocalBranchNames // the remote branches are tracked as local branches in the remote repo

	// initialWorktreeSHAs is only for looking up SHAs that existed at the worktree repo before the first Git Town command was run.
	initialWorktreeSHAs map[string]gitdomain.SHA

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

// InitialBranches provides the branches in this Scenario before the WHEN steps ran.
func (self *ScenarioState) InitialBranches() datatable.DataTable {
	result := datatable.DataTable{}
	result.AddRow("REPOSITORY", "BRANCHES")
	localBranchesJoined := self.initialLocalBranches.Join(", ")
	remoteBranchesJoined := self.initialRemoteBranches.Join(", ")
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

// Reset restores the null value of this ScenarioState.
func (self *ScenarioState) Reset(gitEnv fixture.Fixture) {
	self.fixture = gitEnv
	self.initialLocalBranches = gitdomain.NewLocalBranchNames("main")
	self.initialRemoteBranches = gitdomain.NewLocalBranchNames("main")
	self.initialDevSHAs = map[string]gitdomain.SHA{}
	self.initialOriginSHAs = map[string]gitdomain.SHA{}
	self.initialLineage = datatable.DataTable{Cells: [][]string{{"BRANCH", "PARENT"}}}
	self.initialCurrentBranch = gitdomain.EmptyLocalBranchName()
	self.insideGitRepo = true
	self.runOutput = ""
	self.runExitCode = 0
	self.runExitCodeChecked = false
	self.uncommittedFileName = ""
	self.uncommittedContent = ""
}

// compareExistingCommits compares the commits in the Git environment of the given ScenarioState
// against the given Gherkin table.
func (self *ScenarioState) compareGherkinTable(table *messages.PickleStepArgument_PickleTable) error {
	fields := helpers.TableFields(table)
	commitTable := self.fixture.CommitTable(fields)
	diff, errorCount := commitTable.EqualGherkin(table)
	if errorCount != 0 {
		fmt.Printf("\nERROR! Found %d differences in the existing commits\n\n", errorCount)
		fmt.Println(diff)
		return errors.New("mismatching commits found, see diff above")
	}
	return nil
}
