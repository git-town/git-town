package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
)

// Environments sets up the various environments to test Git operations.
// An environment is a directory in which a Cucumber scenario executes.
// It contains a number of Git repositories.
// The standardized environment in which each Cucumber scenario starts has this setup:
// - the "developer" folder contains a repo that is our workspace
//   (where we run tests in)
// - the "origin" folder contains a repo that is the remote for the developer repo
//   (where pushes from "developer" go to)
// - all repos contain a "main" branch that is configured as Git Town's main branch
//
// Setting up the standardized environment happens a lot (before each scenario).
// To make this process performant,
// a fresh setup is cached in the "memoized" directory
// and copied into the test directory.
type Environments struct {

	// baseDir contains the name of the folder
	// that this class operates in.
	// This folder contains the memoized environment cache
	// as well as the environment used for a current scenario.
	// Once we support concurrency, possibly multiple environments at the same time.
	// This was named REPOSITORY_BASE before.
	baseDir string

	runner *Runner
}

// NewEnvironments creates a new Environments instance
// and prepopulates its environment cache.
func NewEnvironments(runner *Runner) (*Environments, error) {

	// create temp dir
	root, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, errors.Wrap(err, "cannot create temp directory")
	}
	fmt.Println("REPOSITORY_BASE:", root)

	environments := &Environments{baseDir: root, runner: runner}
	err = environments.createMemoizedEnvironment()
	if err != nil {
		return environments, errors.Wrap(err, "Cannot create memoized environment")
	}
	return environments, nil
}

// createMemoizedEnvironment creates a cache for the standardized environment
// that all Cucumber Scenarios start out with,
// including a "main" branch and an "origin" remote.
func (e *Environments) createMemoizedEnvironment() error {

	// Create origin repo and set "main" as default branch
	fmt.Println("creating origin repository")
	repoPath := e.repositoryPath("origin")
	fmt.Println("repository path:", repoPath)
	fmt.Println(0x777)
	err := os.MkdirAll(repoPath, 0777)
	if err != nil {
		return errors.Wrap(err, "cannot create origin directory")
	}
	_, err = e.RunInRepo("origin", "git", "init", "--bare", repoPath)
	if err != nil {
		return errors.Wrap(err, "cannot initialize a bare repo in origin")
	}

	// set "main" as default branch
	_, err = e.RunInRepo("origin", "git", "symbolic-ref", "HEAD", "refs/heads/main")
	if err != nil {
		return errors.Wrap(err, "cannot set 'main' as default branch in origin repo")
	}

	// clone the "developer" repo
	err = e.cloneRepo("origin", "developer")
	if err != nil {
		return errors.Wrap(err, "cannot clone developer repo from origin")
	}

	// Initialize main branch
	err = e.RunManyInRepo("developer", [][]string{
		[]string{"git", "checkout", "--orphan", "main"},
		[]string{"git", "commit", "--allow-empty", "-m", "Initial commit"},
		[]string{"git", "push", "-u", "origin", "main"},
	})
	if err != nil {
		return err
	}

	// memoize environment by saving directory contents
	// FileUtils.cp_r "#{REPOSITORY_BASE}/.", MEMOIZED_REPOSITORY_BASE
	return nil
}

func (e *Environments) cloneRepo(parentName string, childName string) error {

	// clone the repo
	parentPath := e.repositoryPath(parentName)
	childPath := e.repositoryPath(childName)
	_, err := e.runner.Run("git", "clone", parentPath, childPath)
	if err != nil {
		return err
	}

	// configure the repo
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.Chdir(childPath)
	if err != nil {
		return err
	}
	userName := strings.Replace(childName, "_secondary", "", 1)
	err = e.runner.RunMany([][]string{
		[]string{"git", "config", "user.name", userName},
		[]string{"git", "config", "user.email", userName + "@example.com"},
		[]string{"git", "config", "push.default", "simple"},
		[]string{"git", "config", "core.editor", "vim"},
		[]string{"git", "config", "git-town.main-branch-name", "main"},
		[]string{"git", "config", "git-town.perennial-branch-names", ""},
	})
	if err != nil {
		return err
	}
	return os.Chdir(currentDir)
}

func (e *Environments) repositoryPath(repoName string) string {
	return path.Join(e.baseDir, "/", repoName)
}

// runInRepo runs the given command with the given arguments in the given repository.
// It is part of the business logic that the command to run could fail,
// hence errors from the command run aren't signaled in the error return value
// but are contained in the returned Runner instance.
// All other errors (around setting up the call) are returned as the error object
// and must be handled by the called.
func (e *Environments) RunInRepo(repoName string, command string, args ...string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "cannot get current directory")
	}
	repoPath := e.repositoryPath(repoName)
	err = os.Chdir(repoPath)
	if err != nil {
		return "", errors.Wrapf(err, "cannot cd into directory '%s'", repoPath)
	}
	output, err := e.runner.Run(command, args...)
	if err != nil {
		return output, errors.Wrapf(err, "error running %s command\noutput:%s", command, output)
	}
	err = os.Chdir(currentDir)
	if err != nil {
		return "", errors.Wrapf(err, "cannot cd into directory %s", currentDir)
	}
	return output, nil
}

func (e *Environments) RunStringInRepo(repoName, commandText string) (string, error) {
	parts := strings.Fields(commandText)
	command, args := parts[0], parts[1:]
	return e.RunInRepo(repoName, command, args...)
}

// RunManyInRepo runs all given commands in the repo with the given name.
// Failed commands cause abortion of the function with the received error.
func (e *Environments) RunManyInRepo(repoName string, commands [][]string) error {
	for _, commandList := range commands {
		command, args := commandList[0], commandList[1:]
		_, err := e.RunInRepo(repoName, command, args...)
		if err != nil {
			return err
		}
	}
	return nil
}
