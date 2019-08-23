package steps

import "github.com/Originate/git-town/test"

// the GitManager instance to use
var gitManager *test.GitManager

// the GitEnvironment used in the current scenario
var gitEnvironment *test.GitEnvironment

// the result of the last run of Git Town
var lastRunResult test.RunResult
