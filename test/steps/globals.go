package steps

import "github.com/Originate/git-town/test"

// the GitManager instance to use
var gitManager *test.GitManager

// the GitEnvironment used in the current scenario
var gitEnvironment *test.GitEnvironment

// the result of the last run of Git Town
var lastRunOutput string

// the error of the last run of Git Town
var lastRunErr error
