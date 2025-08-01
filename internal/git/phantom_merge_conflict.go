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

// information about a file with merge conflicts, as provided by "git ls-files --unmerged"
type FileConflictQuickInfo struct {
	BaseChange          Option[BlobInfo] // info about the base version of the file (when 3-way merging)
	CurrentBranchChange Option[BlobInfo] // info about the content of the file on the branch where the merge conflict occurs, None == file is deleted here
	IncomingChange      Option[BlobInfo] // info about the content of the file on the branch being merged in, None == file is being deleted here
}

// describes the content of a file blob in Git
type BlobInfo struct {
	FilePath   string        // relative path of the file in the repo
	Permission string        // permissions, in the form "100755"
	SHA        gitdomain.SHA // checksum of the content blob of the file - this is not the commit SHA!
}

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

// Everything Git Town needs to know about a file merge conflict to determine whether this is a phantom merge conflict.
// Includes the FileConflictQuickInfo as well as information that only Git Town knows,
// like how this file looks at the root branch of the stack on which the conflict occurs.
type FileConflictFullInfo struct {
	Current Option[BlobInfo] // info about the file on the current branch
	Parent  Option[BlobInfo] // info about the file on the original parent
	Root    Option[BlobInfo] // info about the file on the root branch
}

// describes a file within an unresolved merge conflict that experiences a phantom merge conflict
type PhantomConflict struct {
	FilePath   string
	Resolution gitdomain.ConflictResolution
}

func DetectPhantomMergeConflicts(conflictInfos []FileConflictFullInfo, parentBranchOpt Option[gitdomain.LocalBranchName], rootBranch gitdomain.LocalBranchName) []PhantomConflict {
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

func EmptyBlobInfo() BlobInfo {
	var result BlobInfo
	return result
}

func ParseLsFilesUnmergedLine(line string) (BlobInfo, UnmergedStage, string, error) {
	// Example text to parse:
	// 100755 ece1e56bf2125e5b114644258872f04bc375ba69 3	file
	permissions, remainder, match := strings.Cut(line, " ")
	if !match {
		return BlobInfo{}, 0, "", fmt.Errorf("cannot read permissions portion from output of \"git ls-files --unmerged\": %q", line)
	}
	shaText, remainder, match := strings.Cut(remainder, " ")
	if !match {
		return BlobInfo{}, 0, "", fmt.Errorf("cannot read SHA portion from output of \"git ls-files --unmerged\": %q", line)
	}
	sha, err := gitdomain.NewSHAErr(shaText)
	if err != nil {
		return BlobInfo{}, 0, "", fmt.Errorf("invalid SHA (%w) in output of \"git ls-files --unmerged\": %q", err, line)
	}
	stageText, remainder, match := strings.Cut(remainder, "\t")
	if !match {
		return BlobInfo{}, 0, "", fmt.Errorf("cannot read stage portion from output of \"git ls-files --unmerged\": %q", line)
	}
	stageInt, err := strconv.Atoi(stageText)
	if err != nil {
		return BlobInfo{}, 0, "", fmt.Errorf("stage portion from output of \"git ls-files --unmerged\" is not a number (%w): %q", err, line)
	}
	stage, err := NewUnmergedStage(stageInt)
	if err != nil {
		return BlobInfo{}, 0, "", fmt.Errorf("unknown stage ID in output of \"git ls-files --unmerged\": %q", line)
	}
	filePath := remainder
	change := BlobInfo{
		FilePath:   filePath,
		Permission: permissions,
		SHA:        sha,
	}
	return change, stage, filePath, nil
}

func ParseLsFilesUnmergedOutput(output string) ([]FileConflictQuickInfo, error) {
	// Example output to parse:
	// 100755 c887ff2255bb9e9440f9456bcf8d310bc8d718d4 2	file
	// 100755 ece1e56bf2125e5b114644258872f04bc375ba69 3	file
	result := []FileConflictQuickInfo{}
	filePathOpt := None[string]()
	baseChange := None[BlobInfo]()
	currentBranchChange := None[BlobInfo]()
	incomingChange := None[BlobInfo]()
	for _, line := range stringslice.Lines(output) {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		change, stage, file, err := ParseLsFilesUnmergedLine(line)
		if err != nil {
			return []FileConflictQuickInfo{}, err
		}
		filePath, hasFilePath := filePathOpt.Get()
		if hasFilePath && file != filePath {
			result = append(result, FileConflictQuickInfo{
				BaseChange:          baseChange,
				CurrentBranchChange: currentBranchChange,
				IncomingChange:      incomingChange,
			})
			filePathOpt = Some(file)
			baseChange = None[BlobInfo]()
			currentBranchChange = None[BlobInfo]()
			incomingChange = None[BlobInfo]()
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
		result = append(result, FileConflictQuickInfo{
			BaseChange:          baseChange,
			CurrentBranchChange: currentBranchChange,
			IncomingChange:      incomingChange,
		})
	}
	return result, nil
}

func ParseLsTreeOutput(output string) (BlobInfo, error) {
	// Example output:
	// 100755 blob ece1e56bf2125e5b114644258872f04bc375ba69	file
	output = strings.TrimSpace(output)
	// skip permissions
	permission, remainder, match := strings.Cut(output, " ")
	if !match {
		return EmptyBlobInfo(), fmt.Errorf("cannot read permissions portion from the output of \"git ls-tree\": %q", output)
	}
	objType, remainder, match := strings.Cut(remainder, " ")
	if !match {
		return EmptyBlobInfo(), fmt.Errorf("cannot read object type from the output of \"git ls-tree\": %q", output)
	}
	if objType != "blob" {
		return EmptyBlobInfo(), fmt.Errorf("unexpected object type (%s) in the output of \"git ls-tree\": %q", objType, output)
	}
	shaText, remainder, match := strings.Cut(remainder, "\t")
	if !match {
		return EmptyBlobInfo(), fmt.Errorf("cannot read SHA from the output of \"git ls-tree\": %q", output)
	}
	sha, err := gitdomain.NewSHAErr(shaText)
	if err != nil {
		return EmptyBlobInfo(), fmt.Errorf("invalid SHA (%s) in the output of \"git ls-tree\": %q", shaText, output)
	}
	blobInfo := BlobInfo{
		FilePath:   remainder,
		Permission: permission,
		SHA:        sha,
	}
	return blobInfo, nil
}
