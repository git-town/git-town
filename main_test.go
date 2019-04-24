package main

/*
Test setup:
- by default, each scenario runs in a directory called "developer"
	that has a "main" branch and a valid Git Town configuration
- at script startup, it creates a memoized repo with that setup
- before each scenario, it copies that memoized repo over into the "developer" repo
*/

import (
	"fmt"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/test"
)

var runner *test.Runner

// repoManager manages the various Git repositories
// that we use for testing the functionality of Git Town.
// type repoManager struct {

// 	// rootDir contains the name of the directory
// 	// that contains the various Git repositories.
// 	// This was named REPOSITORY_BASE before.
// 	rootDir string
// }

func beforeSuite() {
	runner = &test.Runner{}
}

func beforeScenario(interface{}) {
	// copy MEMOIZED_REPOSITORY_BASE to REPOSITORY_BASE
}

// createRepository creates the
// func createOriginRepository() error {
// 	fmt.Println("creating origin repository")
// 	repoPath := repositoryPath("origin")
// 	fmt.Println("repository path:", repoPath)
// 	err := os.MkdirAll(repoPath, 644)
// 	if err != nil {
// 		return err
// 	}
// 	return run("git", "init", "--bare", repoPath)
// }

// func repositoryPath(repoName string) string {
// 	return REPOSITORY_BASE + "/" + repoName
// }

// func runInRepo(repoName string, command string, args ...string) error {
// 	path := repositoryPath(repoName)
// 	// at_path path, &block
// }

// func initializeEnvironment() {

// // Create origin repo and set "main" as default branch
// createRepository("origin")
// runInRepo("origin", "git symbolic-ref HEAD refs/heads/main")

// cloneRepo("origin", "developer")

// // Initialize main branch
// runInRepo("developer", "git checkout --orphan main")
// runInRepo("developer", "git commit --allow-empty -m 'Initial commit'")
// runInRepo("developer", "git push -u origin main")

// // memoize environment by saving directory contents
// // FileUtils.cp_r "#{REPOSITORY_BASE}/.", MEMOIZED_REPOSITORY_BASE
// }

func myWorkspaceIsCurrentlyNotAGitRepository() error {
	// FileUtils.rm_rf '.git'
	return nil
}

func iHaventConfiguredGitTownYet() error {
	// delete_main_branch_configuration
	// delete_perennial_branches_configuration
	return nil
}

func iRun(command string) error {
	runner.RunString(command)
	return nil
}

func itPrints(expected *gherkin.DocString) error {
	if !runner.OutputContains(expected.Content) {
		return fmt.Errorf(`text not found: %s`, expected.Content)
	}
	return nil
}

func itDoesNotPrint(text string) error {
	if runner.OutputContains(text) {
		return fmt.Errorf(`text found: %s`, text)
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	s.BeforeSuite(beforeSuite)
	s.BeforeScenario(beforeScenario)
	s.Step(`^I haven\'t configured Git Town yet$`, iHaventConfiguredGitTownYet)
	s.Step("^my workspace is currently not a Git repository$", myWorkspaceIsCurrentlyNotAGitRepository)
	s.Step(`^I run "([^"]*)"$`, iRun)
	s.Step("^it prints$", itPrints)
	s.Step("^it does not print \"([^\"]*)\"$", itDoesNotPrint)
}
