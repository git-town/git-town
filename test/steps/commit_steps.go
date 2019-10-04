package steps

import (
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

// CommitSteps provides Cucumber step implementations around commits.
func CommitSteps(s *godog.Suite, gtf *GitTownFeature) {
	s.Step(`^the following commit exists in my repository$`, gtf.theFollowingCommitExistsInMyRepository)
}

func (gtf *GitTownFeature) theFollowingCommitExistsInMyRepository(table *gherkin.DataTable) error {
	// user = (who == 'my') ? 'developer' : 'coworker'
	// user += '_secondary' if remote
	// @initial_commits_table = table.clone
	// @original_files = files_in_branches
	// in_repository user do
	fmt.Println("gitEnvironment.DeveloperRepo", gtf.gitEnvironment.DeveloperRepo)
	return gtf.gitEnvironment.DeveloperRepo.CreateCommits(table)
}
