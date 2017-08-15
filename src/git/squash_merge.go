package git

import "github.com/Originate/git-town/src/util"

// CommentOutDefaultSquashCommitMessage comments out the default message for the current squash merge
func CommentOutDefaultSquashCommitMessage() {
	util.GetCommandOutput("sed", "-i", "-e", "s/^/# /g", ".git/SQUASH_MSG")
}
