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
	"os/exec"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

// REPOSITORY_BASE contains the root directory
// in which all Git repos are stored.
// var REPOSITORY_BASE string

// lastRunOutput contains the output of the last command
var lastRunOutput string

func beforeSuite() {
	// REPOSITORY_BASE, err := ioutil.TempDir("", "")
	// if err != nil {
	// 	log.Fatalf("cannot create temp directory", err)
	// }
	// fmt.Println("REPOSITORY_BASE:", REPOSITORY_BASE)
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

// run runs the command consisting of the given elements.
func run(name string, commands ...string) error {
	cmd := exec.Command(name, commands...)
	output, err := cmd.CombinedOutput()
	lastRunOutput = string(output)
	// fmt.Println("output:", lastRunOutput)
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}
	return err
}

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

func iRun(command string) error {
	// NOTE: we split the string by space here, this only works for simple commands without quotes
	parts := strings.Fields(command)
	command, args := parts[0], parts[1:]
	return run(command, args...)
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
	s.Step("^my workspace is currently not a Git repository$", myWorkspaceIsCurrentlyNotAGitRepository)
	s.Step(`^I run "([^"]*)"$`, iRun)
	s.Step("^it prints$", itPrints)
	s.Step("^it does not print \"([^\"]*)\"$", itDoesNotPrint)
}
