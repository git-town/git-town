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
	"io/ioutil"
	"log"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/test"
)

var environments *test.Environments
var runner *test.Runner
var lastRunOutput string
var lastRunError error

func beforeSuite() {
	baseDir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatalf("cannot create base directory: %s", err)
	}
	runner = test.NewRunner(baseDir)
	environments = test.NewEnvironments(baseDir, runner)
	if err != nil {
		log.Fatalf("cannot set up new environment: %s", err)
	}
	err = environments.CreateMemoizedEnvironment()
	if err != nil {
		log.Fatalf("Cannot create memoized environment: %s", err)
	}
}

func beforeScenario(interface{}) {
	// copy MEMOIZED_REPOSITORY_BASE to REPOSITORY_BASE
}

func afterScenario(args interface{}, err error) {
	runner.RemoveTempShellOverrides()
}

func myWorkspaceIsCurrentlyNotAGitRepository() error {
	// FileUtils.rm_rf '.git'
	return nil
}

func iHaveGitInstalled(arg1 string) error {
	err := runner.AddTempShellOverride(
		"git",
		`#!/usr/bin/env bash
		echo "git version 2.6.2"`)
	return err
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

func itPrintsTheError(expected *gherkin.DocString) error {
	if !strings.Contains(lastRunOutput, expected.Content) {
		return fmt.Errorf("text not found: %s\n\nactual text:\n%s", expected.Content, lastRunOutput)
	}
	if lastRunError == nil {
		return fmt.Errorf("expected error")
	}
	return nil
}

func itRunsTheCommands(table *gherkin.DataTable) error {
	commands := test.CommandsInOutput(lastRunOutput)
	return test.AssertStringSliceMatchesTable(commands, table)
}

func itRunsNoCommands() error {
	commands := test.CommandsInOutput(lastRunOutput)
	if len(commands) > 0 {
		for _, command := range commands {
			fmt.Println(command)
		}
		return fmt.Errorf("expected no commands but found %d commands", len(commands))
	}
	return nil
}

func theFollowingCommitExistsInMyRepository(table *gherkin.DataTable) error {
	return godog.ErrPending
}

// nolint:deadcode
func FeatureContext(s *godog.Suite) {
	s.BeforeSuite(beforeSuite)
	s.BeforeScenario(beforeScenario)
	s.AfterScenario(afterScenario)
	s.Step(`^I haven\'t configured Git Town yet$`, iHaventConfiguredGitTownYet)
	s.Step("^my workspace is currently not a Git repository$", myWorkspaceIsCurrentlyNotAGitRepository)
	s.Step(`^I run "([^"]*)"$`, iRun)
	s.Step("^it prints$", itPrints)
	s.Step("^it does not print \"([^\"]*)\"$", itDoesNotPrint)
	s.Step(`^it prints the error:$`, itPrintsTheError)
	s.Step(`^I have Git "([^"]*)" installed$`, iHaveGitInstalled)
	s.Step(`^it runs the commands$`, itRunsTheCommands)
	s.Step(`^it runs no commands$`, itRunsNoCommands)
	s.Step(`^the following commit exists in my repository$`, theFollowingCommitExistsInMyRepository)
}
