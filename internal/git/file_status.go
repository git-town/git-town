package git

import (
	"fmt"
	"slices"
	"strings"

	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// FileStatus contains information about each file reported by "git status -z".
type FileStatus struct {
	OriginalPath Option[string] // if the file was renamed or copied, the old path of the file
	Path         string         // the path to the file
	ShortStatus  string         // a two-letter status code, explained at https://git-scm.com/docs/git-status#_short_format
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
		originalPath := None[string]()
		if strings.Contains(status, "R") || strings.Contains(status, "C") {
			i++
			if i >= len(entries) {
				return nil, fmt.Errorf(messages.InvalidStatusOutput, entry)
			}
			originalPathString := entries[i]
			if originalPathString == "" {
				return nil, fmt.Errorf(messages.InvalidStatusOutput, entry)
			}
			originalPath = Some(originalPathString)
		}
		result = append(result, FileStatus{
			OriginalPath: originalPath,
			Path:         path,
			ShortStatus:  status,
		})
	}
	return result, nil
}
