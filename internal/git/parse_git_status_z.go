package git

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v19/internal/messages"
	. "github.com/git-town/git-town/v19/pkg/prelude"
)

// information about each file reported by "git status -z"
type FileStatus struct {
	Status      string         // a two-letter status code, explained at https://git-scm.com/docs/git-status#_short_format
	Path        string         // the path to the file
	RenamedFrom Option[string] // if the file was renamed, the old path of the file
}

// ParseGitStatusZ parses the output of "git status -z" into a slice of FileStatus.
func ParseGitStatusZ(output string) ([]FileStatus, error) {
	var result []FileStatus
	entries := strings.Split(output, "\x00")
	for i := 0; i < len(entries); i++ {
		entry := entries[i]
		if entry == "" {
			continue
		}
		if len(entry) < 4 {
			return nil, fmt.Errorf(messages.InvalidStatusOutput, entry)
		}
		status := entry[:2]
		if entry[2] != ' ' {
			return nil, fmt.Errorf(messages.InvalidStatusOutput, entry)
		}
		path := entry[3:]
		renamedFrom := None[string]()
		if strings.Contains(status, "R") {
			i++
			if i >= len(entries) {
				return nil, fmt.Errorf(messages.InvalidStatusOutput, entry)
			}
			renamedFromString := entries[i]
			if renamedFromString == "" {
				return nil, fmt.Errorf(messages.InvalidStatusOutput, entry)
			}
			renamedFrom = Some(renamedFromString)
		}
		result = append(result, FileStatus{
			Status:      status,
			Path:        path,
			RenamedFrom: renamedFrom,
		})
	}
	return result, nil
}
