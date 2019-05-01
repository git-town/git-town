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

	"github.com/dchest/uniuri"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/git-town/test"
	"github.com/iancoleman/strcase"
)

// the GitManager instance to use
var gitManager *test.GitManager

// the GitEnvironment used in the current scenario
var gitEnvironment *test.GitEnvironment

// the result of the last run of Git Town
var lastRunResult test.RunResult

func beforeSuite() {

	// create the directory to put the GitEnvironments ino
	baseDir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatalf("cannot create base directory: %s", err)
	}

	// create the GitManager
	gitManager = test.NewGitManager(baseDir)

	// create the memoized environment
	err = gitManager.CreateMemoizedEnvironment()
	if err != nil {
		log.Fatalf("Cannot create memoized environment: %s", err)
	}
}

func beforeScenario(args interface{}) {

	// create a GitEnvironment for the scenario
	environmentName := strcase.ToSnake(scenarioName(args)) + "_" + string(uniuri.NewLen(10))
	var err error
	gitEnvironment, err = gitManager.CreateScenarioEnvironment(environmentName)
	if err != nil {
		log.Fatalf("cannot create environment for scenario '%s': %s", environmentName, err)
	}
	fmt.Println("GIT ENVIRONMENT", gitEnvironment)
}

func afterScenario(args interface{}, err error) {
	// TODO: delete scenario environment
}

// scenarioName returns the name of the given Scenario or ScenarioOutline
func scenarioName(args interface{}) string {
	scenario, ok := args.(*gherkin.Scenario)
	if ok {
		return scenario.Name
	}
	scenarioOutline, ok := args.(*gherkin.ScenarioOutline)
	if ok {
		return scenarioOutline.Name
	}
	panic("unknown type")
}

func myWorkspaceIsCurrentlyNotAGitRepository() error {
	// FileUtils.rm_rf '.git'
	return nil
}

func iHaveGitInstalled(arg1 string) error {
	err := gitEnvironment.DeveloperRepo.AddTempShellOverride(
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
	lastRunResult = gitEnvironment.DeveloperRepo.RunString(command)
	return nil
}

func itPrints(expected *gherkin.DocString) error {
	if !strings.Contains(lastRunResult.Output, expected.Content) {
		return fmt.Errorf(`text not found: %s`, expected.Content)
	}
	return nil
}

func itDoesNotPrint(text string) error {
	if strings.Contains(lastRunResult.Output, text) {
		return fmt.Errorf(`text found: %s`, text)
	}
	return nil
}

func itPrintsTheError(expected *gherkin.DocString) error {
	if !strings.Contains(lastRunResult.Output, expected.Content) {
		return fmt.Errorf("text not found: %s\n\nactual text:\n%s", expected.Content, lastRunResult.Output)
	}
	if lastRunResult.Err == nil {
		return fmt.Errorf("expected error")
	}
	return nil
}

func itRunsTheCommands(table *gherkin.DataTable) error {
	commands := test.CommandsInOutput(lastRunResult.Output)
	return test.AssertStringSliceMatchesTable(commands, table)
}

func itRunsNoCommands() error {
	commands := test.CommandsInOutput(lastRunResult.Output)
	if len(commands) > 0 {
		for _, command := range commands {
			fmt.Println(command)
		}
		return fmt.Errorf("expected no commands but found %d commands", len(commands))
	}
	return nil
}

func theFollowingCommitExistsInMyRepository(table *gherkin.DataTable) error {
	// user = (who == 'my') ? 'developer' : 'coworker'
	// user += '_secondary' if remote
	// @initial_commits_table = table.clone
	// @original_files = files_in_branches
	// in_repository user do
	fmt.Println("gitEnvironment.DeveloperRepo", gitEnvironment.DeveloperRepo)
	return gitEnvironment.DeveloperRepo.CreateCommits(table)
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
