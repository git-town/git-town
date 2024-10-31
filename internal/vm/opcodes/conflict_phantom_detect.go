package opcodes

import (
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

type ConflictPhantomDetect struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

/*
	"git ls-files --unmerged" shows the status of all conflicting files.

	Example output:

	100755 c887ff2255bb9e9440f9456bcf8d310bc8d718d4 2	file
	100755 ece1e56bf2125e5b114644258872f04bc375ba69 3	file

	The first column is the file permissions. They must match in phantom merge conflicts.
	The second column is the SHA1 of the file content (blob).
	The third column is the stage number.
	The fourth column is the path of the conflicting file.

	Stage 1 (not in the example output) is the common ancestor (base version).
	Stage 2 (c887ff) is the the local file (on the current branch)
	Stage 3 (ece1e5) is the version being merged in (from the parent branch)

	To see the SHA1 of the content blob of a file on a particular branch:

	git ls-tree <branch> <file-path>
	git ls-tree main file

	Example output:

	100755 blob ece1e56bf2125e5b114644258872f04bc375ba69	file

	The file is identical if the SHA1 of the blobs is identical.
	In this case, the SHA1 of the incoming file version (from stage 3) is identical to the SHA1 of that file on the main branch: ece1e5.
	This means it's a phantom merge conflict.
*/

func (self *ConflictPhantomDetect) Run(args shared.RunArgs) error {
	unmergedFiles, err := args.Git.UnmergedFiles(args.Backend)
	if err != nil {
		return err
	}
	phantomMergeConflicts, err := args.Git.DetectPhantomMergeConflicts(args.Backend, unmergedFiles, args.Config.Value.ValidatedConfigData.MainBranch)
	newOpcodes := make([]shared.Opcode, len(phantomMergeConflicts)+1)
	for p, phantomMergeConflict := range phantomMergeConflicts {
		newOpcodes[p] = &ConflictPhantomResolve{
			FilePath: phantomMergeConflict.FilePath,
		}
	}
	newOpcodes[len(phantomMergeConflicts)] = &ConflictPhantomFinalize{}
	return nil
}
