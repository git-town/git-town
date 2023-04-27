package commands

import (
	"fmt"
	"strings"
)

// Tags provides a list of the tags in this repository.
func Tags(shell Shell) ([]string, error) {
	output, err := shell.Run("git", "tag")
	if err != nil {
		return []string{}, fmt.Errorf("cannot determine tags in repo %q: %w", shell.Dir(), err)
	}
	result := []string{}
	for _, line := range strings.Split(output, "\n") {
		result = append(result, strings.TrimSpace(line))
	}
	return result, err
}
