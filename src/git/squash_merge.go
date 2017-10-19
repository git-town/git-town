package git

import (
	"io/ioutil"
	"regexp"

	"github.com/Originate/exit"
)

var squashMessageFile = ".git/SQUASH_MSG"

// CommentOutSquashCommitMessage comments out the message for the current squash merge
// Adds the given prefix with the newline if provided
func CommentOutSquashCommitMessage(prefix string) {
	contentBytes, err := ioutil.ReadFile(squashMessageFile)
	exit.If(err)
	content := string(contentBytes)
	if prefix != "" {
		content = prefix + "\n" + content
	}
	content = regexp.MustCompile("(?m)^").ReplaceAllString(content, "# ")
	err = ioutil.WriteFile(squashMessageFile, []byte(content), 0644)
	exit.If(err)
}
