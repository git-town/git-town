package cucumber

import (
	"errors"
	"fmt"

	"github.com/cucumber/godog"

	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/test/datatable"
	"github.com/git-town/git-town/v20/internal/test/fixture"
	"github.com/git-town/git-town/v20/internal/test/helpers"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

// ScenarioState constains the state that is shared by all steps within a scenario.
type ScenarioState struct {
	// commits before the end-to-end test executed the most recent subshell command
	beforeRunDevSHAs Option[map[string]gitdomain.SHA]

	// commits at the origin remote before the end-to-end test executed the most recent subshell command
	beforeRunOriginSHAs Option[map[string]gitdomain.SHA]

	// commits at the worktree-repo before the end-to-end test executed the most recent shell command
	beforeRunWorktreeSHAs Option[map[string]gitdomain.SHA]

	// the Fixture used in the current scenario
	fixture fixture.Fixture

	// the local and remote branches before the end-to-end test executed the first subshell command
	initialBranches Option[datatable.DataTable]

	// the commits in this Git environment before the end-to-end test executed the first subshell command
	initialCommits Option[datatable.DataTable]

	// the branch that was checked out before the end-to-end test executed the first subshell command
	initialCurrentBranch Option[gitdomain.LocalBranchName]

	// initialDevSHAs is only for looking up SHAs that existed at the developer repo before the first Git Town command ran.
	// It's not a source of truth for which branches existed at that time
	// because it might contain non-existing remote branches or miss existing remote branches.
	// An example is when origin removes a branch. initialDevSHAs will still list it
	// because the developer workspace hasn't fetched updates yet.
	initialDevSHAs Option[map[string]gitdomain.SHA]

	// the lineage before the end-to-end test executed the first subshell command
	initialLineage Option[datatable.DataTable]

	// commits that existed at the origin repo before the end-to-end test executed the first subshell command
	initialOriginSHAs Option[map[string]gitdomain.SHA]

	// the Git tags before the end-to-end test executed the first subshell command
	initialTags Option[datatable.DataTable]

	// commits at the worktree repo before the end-to-end test executed the first subshell command
	initialWorktreeSHAs Option[map[string]gitdomain.SHA]

	// whether the developer workspace contains a Git repository
	insideGitRepo bool

	// the error of the last run of Git Town
	runExitCode Option[int]

	// indicates whether the scenario has verified the error
	runExitCodeChecked bool

	// the output of the last run of Git Town
	runOutput Option[string]

	// content of the uncommitted file in the workspace
	uncommittedContent Option[string]

	// name of the uncommitted file in the workspace
	uncommittedFileName Option[string]
}

func (self *ScenarioState) CaptureState() {
	if self.initialCommits.IsNone() && self.insideGitRepo && self.fixture.SubmoduleRepo.IsNone() {
		currentCommits := self.fixture.CommitTable([]string{"BRANCH", "LOCATION", "MESSAGE", "FILE NAME", "FILE CONTENT"})
		self.initialCommits = Some(currentCommits)
	}
	if self.initialBranches.IsNone() && self.insideGitRepo {
		branches := self.fixture.Branches()
		self.initialBranches = Some(branches)
	}
	if self.initialLineage.IsNone() && self.insideGitRepo {
		lineage := self.fixture.DevRepo.GetOrPanic().LineageTable()
		self.initialLineage = Some(lineage)
	}
	if self.initialTags.IsNone() && self.insideGitRepo {
		tags := self.fixture.TagTable()
		self.initialTags = Some(tags)
	}
}

// compareExistingCommits compares the commits in the Git environment of the given ScenarioState
// against the given Gherkin table.
func (self *ScenarioState) compareGherkinTable(table *godog.Table) error {
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
