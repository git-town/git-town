package git

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/gohacks/stringslice"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// information about a file with merge conflicts
type UnmergedFile struct {
	BaseChange          Option[LsFilesUnmergedChange]
	CurrentBranchChange LsFilesUnmergedChange
	FilePath            string
	IncomingChange      LsFilesUnmergedChange
}

func (self UnmergedFile) HasDifferentPermissions() bool {
	if self.CurrentBranchChange.Permission != self.IncomingChange.Permission {
		return true
	}
	if baseChange, hasBaseChange := self.BaseChange.Get(); hasBaseChange {
		if baseChange.Permission != self.IncomingChange.Permission {
			return true
		}
	}
	return false
}

type LsFilesUnmergedChange struct {
	Permission string
	SHA        gitdomain.SHA
}

type LsFilesUnmergedStage int

const (
	LsFilesUnmergedStageBase          LsFilesUnmergedStage = 1
	LsFilesUnmergedStageCurrentBranch LsFilesUnmergedStage = 2
	LsFilesUnmergedStageIncoming      LsFilesUnmergedStage = 3
)

var LsFilesUnmergedStages = []LsFilesUnmergedStage{ //nolint:gochecknoglobals
	LsFilesUnmergedStageBase,
	LsFilesUnmergedStageCurrentBranch,
	LsFilesUnmergedStageIncoming,
}

func ParseLsFilesUnmergedLine(line string) (LsFilesUnmergedChange, LsFilesUnmergedStage, string, error) {
	// Example output to parse:
	// 100755 c887ff2255bb9e9440f9456bcf8d310bc8d718d4 2	file
	// 100755 ece1e56bf2125e5b114644258872f04bc375ba69 3	file
	permissions, remainder, match := strings.Cut(line, " ")
	if !match {
		return LsFilesUnmergedChange{}, 0, "", fmt.Errorf("cannot read permissions portion from output of \"git ls-files --unmerged\": %q", line)
	}
	shaText, remainder, match := strings.Cut(remainder, " ")
	if !match {
		return LsFilesUnmergedChange{}, 0, "", fmt.Errorf("cannot read SHA portion from output of \"git ls-files --unmerged\": %q", line)
	}
	sha, err := gitdomain.NewSHAErr(shaText)
	if err != nil {
		return LsFilesUnmergedChange{}, 0, "", fmt.Errorf("invalid SHA (%w) in output of \"git ls-files --unmerged\": %q", err, line)
	}
	stageText, remainder, match := strings.Cut(remainder, "\t")
	if !match {
		return LsFilesUnmergedChange{}, 0, "", fmt.Errorf("cannot read stage portion from output of \"git ls-files --unmerged\": %q", line)
	}
	stageInt, err := strconv.Atoi(stageText)
	if err != nil {
		return LsFilesUnmergedChange{}, 0, "", fmt.Errorf("stage portion from output of \"git ls-files --unmerged\" is not a number (%w): %q", err, line)
	}
	stage, err := NewLsFilesUnmergedStage(stageInt)
	if err != nil {
		return LsFilesUnmergedChange{}, 0, "", fmt.Errorf("unknown stage ID in output of \"git ls-files --unmerged\": %q", line)
	}
	filePath := remainder
	change := LsFilesUnmergedChange{
		Permission: permissions,
		SHA:        sha,
	}
	return change, stage, filePath, nil
}

func ParseLsFilesUnmergedOutput(output string) ([]UnmergedFile, error) {
	result := []UnmergedFile{}
	filePath := ""
	baseChangeOpt := None[LsFilesUnmergedChange]()
	currentBranchChangeOpt := None[LsFilesUnmergedChange]()
	incomingChangeOpt := None[LsFilesUnmergedChange]()
	for _, line := range stringslice.Lines(output) {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		change, stage, file, err := ParseLsFilesUnmergedLine(line)
		if err != nil {
			return []UnmergedFile{}, err
		}
		if file != filePath {
			currentBranchChange, hasCurrentBranchChange := currentBranchChangeOpt.Get()
			incomingChange, hasIncomingChange := incomingChangeOpt.Get()
			if hasCurrentBranchChange && hasIncomingChange {
				unmergedFile := UnmergedFile{
					BaseChange:          baseChangeOpt,
					CurrentBranchChange: currentBranchChange,
					FilePath:            filePath,
					IncomingChange:      incomingChange,
				}
				result = append(result, unmergedFile)
			}
			filePath = file
			baseChangeOpt = None[LsFilesUnmergedChange]()
			currentBranchChangeOpt = None[LsFilesUnmergedChange]()
			incomingChangeOpt = None[LsFilesUnmergedChange]()
		}
		switch stage {
		case LsFilesUnmergedStageBase:
			baseChangeOpt = Some(change)
		case LsFilesUnmergedStageCurrentBranch:
			currentBranchChangeOpt = Some(change)
		case LsFilesUnmergedStageIncoming:
			incomingChangeOpt = Some(change)
		}
	}
	if len(filePath) > 0 {
		currentBranchChange, hasCurrentBranchChange := currentBranchChangeOpt.Get()
		incomingChange, hasIncomingChange := incomingChangeOpt.Get()
		if hasCurrentBranchChange && hasIncomingChange {
			unmergedFile := UnmergedFile{
				BaseChange:          baseChangeOpt,
				CurrentBranchChange: currentBranchChange,
				FilePath:            filePath,
				IncomingChange:      incomingChange,
			}
			result = append(result, unmergedFile)
		}
	}
	return result, nil
}

func ParseLsTreeOutput(output string) (gitdomain.SHA, error) {
	// Example output:
	// 100755 blob ece1e56bf2125e5b114644258872f04bc375ba69	file
	output = strings.TrimSpace(output)
	// skip permissions
	_, remainder, match := strings.Cut(output, " ")
	if !match {
		return "", fmt.Errorf("cannot read permissions portion from the output of \"git ls-tree\": %q", output)
	}
	objType, remainder, match := strings.Cut(remainder, " ")
	if !match {
		return "", fmt.Errorf("cannot read object type from the output of \"git ls-tree\": %q", output)
	}
	if objType != "blob" {
		return "", fmt.Errorf("unexpected object type (%s) in the output of \"git ls-tree\": %q", objType, output)
	}
	shaText, _, match := strings.Cut(remainder, "\t")
	if !match {
		return "", fmt.Errorf("cannot read SHA from the output of \"git ls-tree\": %q", output)
	}
	sha, err := gitdomain.NewSHAErr(shaText)
	if err != nil {
		return "", fmt.Errorf("invalid SHA (%s) in the output of \"git ls-tree\": %q", shaText, output)
	}
	return sha, nil
}
