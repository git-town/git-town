package git

import . "github.com/git-town/git-town/v21/pkg/prelude"

// information about a file with merge conflicts, as provided by "git ls-files --unmerged"
type FileConflict struct {
	BaseChange          Option[Blob] // info about the base version of the file (when 3-way merging)
	CurrentBranchChange Option[Blob] // info about the content of the file on the branch where the merge conflict occurs, None == file is deleted here
	IncomingChange      Option[Blob] // info about the content of the file on the branch being merged in, None == file is being deleted here
}
