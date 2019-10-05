package steps

import (
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

// CommitSteps defines Cucumber step implementations around commits.
func CommitSteps(s *godog.Suite, state *FeatureState) {
	s.Step(`^the following commit exists in my repository$`, state.theFollowingCommitExistsInMyRepository)
}

func (state *FeatureState) theFollowingCommitExistsInMyRepository(table *gherkin.DataTable) error {
	// user = (who == 'my') ? 'developer' : 'coworker'
	// user += '_secondary' if remote
	// @initial_commits_table = table.clone
	// @original_files = files_in_branches
	// in_repository user do
	fmt.Println("gitEnvironment.DeveloperRepo", state.gitEnvironment.DeveloperRepo)
	return state.gitEnvironment.DeveloperRepo.CreateCommits(table)
}
