package git

import (
	"io/ioutil"
	"strings"

	"github.com/Originate/git-town/src/exit"
)

var squashMessageFile = ".git/SQUASH_MSG"

// CommentOutSquashCommitMessage comments out the message for the current squash merge
// Adds the given prefix
func CommentOutSquashCommitMessage(prefix string) {
	contentBytes, err := ioutil.ReadFile(squashMessageFile)
	exit.On(err)
	content := string(contentBytes)
	if prefix != "" {
		content = prefix + "\n" + content
	}
	lines := strings.Split(content, "\n")
	for i := range lines {
		lines[i] = "# " + lines[i]
	}
	ioutil.WriteFile(squashMessageFile, []byte(strings.Join(lines, "\n")), 0644)
}
