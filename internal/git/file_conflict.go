package git

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// FileConflict contains information about a file with conflicts, as provided by "git ls-files --unmerged".
type FileConflict struct {
	BaseChange          Option[Blob] // info about the base version of the file (when 3-way merging)
	CurrentBranchChange Option[Blob] // info about the content of the file on the branch where the merge conflict occurs, None == file is deleted here
	IncomingChange      Option[Blob] // info about the content of the file on the branch being merged in, None == file is being deleted here
}

// Debug prints debug information.
func (self FileConflict) Debug(querier subshelldomain.Querier) {
	base, hasBase := self.BaseChange.Get()
	current, hasCurrent := self.CurrentBranchChange.Get()
	incoming, hasIncoming := self.IncomingChange.Get()
	fmt.Print("BASE CHANGE: ")
	if hasBase {
		base.Debug(querier)
	} else {
		fmt.Println("(none)")
	}
	fmt.Print("CURRENT CHANGE: ")
	if hasCurrent {
		current.Debug(querier)
	} else {
		fmt.Println("(none)")
	}
	fmt.Print("INCOMING CHANGE: ")
	if hasIncoming {
		incoming.Debug(querier)
	} else {
		fmt.Println("(none)")
	}
}

func ParseLsFilesUnmergedLine(line string) (Blob, UnmergedStage, string, error) {
	// Example text to parse:
	// 100755 ece1e56bf2125e5b114644258872f04bc375ba69 3  file
	permissions, remainder, match := strings.Cut(line, " ")
	if !match {
		return Blob{}, 0, "", fmt.Errorf("cannot read permissions portion from output of \"git ls-files --unmerged\": %q", line)
	}
	shaText, remainder, match := strings.Cut(remainder, " ")
	if !match {
		return Blob{}, 0, "", fmt.Errorf("cannot read SHA portion from output of \"git ls-files --unmerged\": %q", line)
	}
	sha, err := gitdomain.NewSHAErr(shaText)
	if err != nil {
		return Blob{}, 0, "", fmt.Errorf("invalid SHA (%w) in output of \"git ls-files --unmerged\": %q", err, line)
	}
	stageText, remainder, match := strings.Cut(remainder, "\t")
	if !match {
		return Blob{}, 0, "", fmt.Errorf("cannot read stage portion from output of \"git ls-files --unmerged\": %q", line)
	}
	stageInt, err := strconv.Atoi(stageText)
	if err != nil {
		return Blob{}, 0, "", fmt.Errorf("stage portion from output of \"git ls-files --unmerged\" is not a number (%w): %q", err, line)
	}
	stage, err := NewUnmergedStage(stageInt)
	if err != nil {
		return Blob{}, 0, "", fmt.Errorf("unknown stage ID in output of \"git ls-files --unmerged\": %q", line)
	}
	filePath := remainder
	change := Blob{
		FilePath:   filePath,
		Permission: permissions,
		SHA:        sha,
	}
	return change, stage, filePath, nil
}

func ParseLsFilesUnmergedOutput(output string) (FileConflicts, error) {
	// Example output to parse:
	// 100755 c887ff2255bb9e9440f9456bcf8d310bc8d718d4 2	file
	// 100755 ece1e56bf2125e5b114644258872f04bc375ba69 3	file
	result := FileConflicts{}
	filePathOpt := None[string]()
	baseChange := None[Blob]()
	currentBranchChange := None[Blob]()
	incomingChange := None[Blob]()
	for _, line := range stringslice.Lines(output) {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		change, stage, file, err := ParseLsFilesUnmergedLine(line)
		if err != nil {
			return FileConflicts{}, err
		}
		filePath, hasFilePath := filePathOpt.Get()
		if hasFilePath && file != filePath {
			result = append(result, FileConflict{
				BaseChange:          baseChange,
				CurrentBranchChange: currentBranchChange,
				IncomingChange:      incomingChange,
			})
			filePathOpt = Some(file)
			baseChange = None[Blob]()
			currentBranchChange = None[Blob]()
			incomingChange = None[Blob]()
		}
		switch stage {
		case UnmergedStageBase:
			baseChange = Some(change)
		case UnmergedStageCurrentBranch:
			currentBranchChange = Some(change)
		case UnmergedStageIncoming:
			incomingChange = Some(change)
		}
	}
	if baseChange.IsSome() || currentBranchChange.IsSome() || incomingChange.IsSome() {
		result = append(result, FileConflict{
			BaseChange:          baseChange,
			CurrentBranchChange: currentBranchChange,
			IncomingChange:      incomingChange,
		})
	}
	return result, nil
}
