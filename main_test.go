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
	"log"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/test"
)

var environments *test.Environments
var lastRunOutput string
var lastRunError error

func beforeSuite() {
	var err error
	environments, err = test.NewEnvironments()
	if err != nil {
		log.Fatalf("cannot set up new environment: %s", err)
	}
}

func beforeScenario(interface{}) {
	// copy MEMOIZED_REPOSITORY_BASE to REPOSITORY_BASE
}

// func runInRepo(repoName string, command string, args ...string) error {
// 	path := repositoryPath(repoName)
// 	// at_path path, &block
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
	lastRunOutput, lastRunError = environments.RunStringInRepo("developer", command)
	return nil
}

func itPrints(expected *gherkin.DocString) error {
	if !strings.Contains(lastRunOutput, expected.Content) {
		return fmt.Errorf(`text not found: %s`, expected.Content)
	}
	return nil
}

func itDoesNotPrint(text string) error {
	if strings.Contains(lastRunOutput, text) {
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
