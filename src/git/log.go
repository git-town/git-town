package git

import "github.com/Originate/git-town/src/util"

// GetLastCommitMessage returns the commit message for the last commit
func GetLastCommitMessage() string {
	return util.GetCommandOutput("git", "log", "-1", "--format=%B")
}
