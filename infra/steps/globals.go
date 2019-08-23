package steps

import "github.com/Originate/git-town/infra"

// the GitManager instance to use
var gitManager *infra.GitManager

// the GitEnvironment used in the current scenario
var gitEnvironment *infra.GitEnvironment

// the result of the last run of Git Town
var lastRunResult infra.RunResult
