package git

import (
	"io/ioutil"
	"regexp"

	"github.com/pkg/errors"
)

var squashMessageFile = ".git/SQUASH_MSG"

// CommentOutSquashCommitMessage comments out the message for the current squash merge
// Adds the given prefix with the newline if provided
func CommentOutSquashCommitMessage(prefix string) error {
	contentBytes, err := ioutil.ReadFile(squashMessageFile)
	if err != nil {
		return errors.Wrapf(err, "cannot read squash message file %q", squashMessageFile)
	}
	content := string(contentBytes)
	if prefix != "" {
		content = prefix + "\n" + content
	}
	content = regexp.MustCompile("(?m)^").ReplaceAllString(content, "# ")
	return ioutil.WriteFile(squashMessageFile, []byte(content), 0644)
}
