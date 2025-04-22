package git

import (
	"fmt"
	"slices"
	"strings"

	"github.com/git-town/git-town/v19/internal/messages"
	. "github.com/git-town/git-town/v19/pkg/prelude"
)

// information about each file reported by "git status -z"
type FileStatus struct {
	Path        string         // the path to the file
	RenamedFrom Option[string] // if the file was renamed, the old path of the file
	ShortStatus string         // a two-letter status code, explained at https://git-scm.com/docs/git-status#_short_format
}

func FileStatusIsUnmerged(status FileStatus) bool {
	return slices.Contains([]string{"DD", "AU", "UD", "UA", "DU", "AA", "UU"}, status.ShortStatus)
}

func FileStatusIsUntracked(status FileStatus) bool {
	return status.ShortStatus == "??"
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
			Path:        path,
			RenamedFrom: renamedFrom,
			ShortStatus: status,
		})
	}
	return result, nil
}
