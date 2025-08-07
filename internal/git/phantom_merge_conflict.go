package git

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// describes the roles that a file can play in a merge conflict
type UnmergedStage int

const (
	UnmergedStageBase          UnmergedStage = 1 // the base version in a 3-way merge
	UnmergedStageCurrentBranch UnmergedStage = 2 // the file on the branch on which the merge conflict happens
	UnmergedStageIncoming      UnmergedStage = 3 // the file on the branch getting merged in
)

// all possile UnmergedStages instances
var UnmergedStages = []UnmergedStage{
	UnmergedStageBase,
	UnmergedStageCurrentBranch,
	UnmergedStageIncoming,
}

func DetectPhantomMergeConflicts(conflictInfos []MergeConflict, parentBranchOpt Option[gitdomain.LocalBranchName], rootBranch gitdomain.LocalBranchName) []PhantomConflict {
	parentBranch, hasParentBranch := parentBranchOpt.Get()
	if !hasParentBranch || parentBranch == rootBranch {
		// branches that don't have a parent or whose parent is the root branch cannot have phantom merge conflicts
		return []PhantomConflict{}
	}
	result := []PhantomConflict{}
	for _, conflictInfo := range conflictInfos {
		initialParentInfo, hasInitialParentInfo := conflictInfo.Parent.Get()
		currentInfo, hasCurrentInfo := conflictInfo.Current.Get()
		if !hasInitialParentInfo || !hasCurrentInfo || currentInfo.Permission != initialParentInfo.Permission {
			continue
		}
		if reflect.DeepEqual(conflictInfo.Root, conflictInfo.Parent) {
			// root and parent have the exact same version of the file --> this is a phantom merge conflict
			result = append(result, PhantomConflict{
				FilePath:   currentInfo.FilePath,
				Resolution: gitdomain.ConflictResolutionOurs,
			})
		}
	}
	return result
}

func ParseLsFilesUnmergedLine(line string) (Blob, UnmergedStage, string, error) {
	// Example text to parse:
	// 100755 ece1e56bf2125e5b114644258872f04bc375ba69 3	file
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

func ParseLsFilesUnmergedOutput(output string) ([]FileConflict, error) {
	// Example output to parse:
	// 100755 c887ff2255bb9e9440f9456bcf8d310bc8d718d4 2	file
	// 100755 ece1e56bf2125e5b114644258872f04bc375ba69 3	file
	result := []FileConflict{}
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
			return []FileConflict{}, err
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

func ParseLsTreeOutput(output string) (Blob, error) {
	// Example output:
	// 100755 blob ece1e56bf2125e5b114644258872f04bc375ba69	file
	output = strings.TrimSpace(output)
	// skip permissions
	permission, remainder, match := strings.Cut(output, " ")
	if !match {
		return EmptyBlob(), fmt.Errorf("cannot read permissions portion from the output of \"git ls-tree\": %q", output)
	}
	objType, remainder, match := strings.Cut(remainder, " ")
	if !match {
		return EmptyBlob(), fmt.Errorf("cannot read object type from the output of \"git ls-tree\": %q", output)
	}
	if objType != "blob" {
		return EmptyBlob(), fmt.Errorf("unexpected object type (%s) in the output of \"git ls-tree\": %q", objType, output)
	}
	shaText, remainder, match := strings.Cut(remainder, "\t")
	if !match {
		return EmptyBlob(), fmt.Errorf("cannot read SHA from the output of \"git ls-tree\": %q", output)
	}
	sha, err := gitdomain.NewSHAErr(shaText)
	if err != nil {
		return EmptyBlob(), fmt.Errorf("invalid SHA (%s) in the output of \"git ls-tree\": %q", shaText, output)
	}
	blobInfo := Blob{
		FilePath:   remainder,
		Permission: permission,
		SHA:        sha,
	}
	return blobInfo, nil
}
