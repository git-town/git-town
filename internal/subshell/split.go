package subshell

import (
	"fmt"

	"github.com/kballard/go-shellquote"
)

func Split(command string) (executable string, args []string, err error) {
	words, err := shellquote.Split(command)
	if err != nil {
		return "", []string{}, fmt.Errorf("cannot split shell call (%s): %w", command, err)
	}
	if len(words) == 0 {
		return "", []string{}, fmt.Errorf("shell call %q seems to be empty", command)
	}
	return words[0], words[1:], nil
}
