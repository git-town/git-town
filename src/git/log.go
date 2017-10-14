package git

import "github.com/Originate/git-town/src/runner"

// GetLastCommitMessage returns the commit message for the last commit
func GetLastCommitMessage() string {
	return runner.New("git", "log", "-1", "--format=%B").Output()
}
