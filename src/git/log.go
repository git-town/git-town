package git

import "github.com/Originate/git-town/src/command"

// GetLastCommitMessage returns the commit message for the last commit
func GetLastCommitMessage() string {
	return command.Run("git", "log", "-1", "--format=%B").Output()
}
