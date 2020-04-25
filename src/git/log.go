package git

import "github.com/git-town/git-town/src/command"

// GetLastCommitMessage returns the commit message for the last commit
func GetLastCommitMessage() string {
	return command.MustRun("git", "log", "-1", "--format=%B").OutputSanitized()
}
