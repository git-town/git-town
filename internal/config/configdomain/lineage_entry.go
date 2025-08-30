package configdomain

import "github.com/git-town/git-town/v21/internal/git/gitdomain"

type LineageEntry struct {
	Child  gitdomain.LocalBranchName
	Parent gitdomain.LocalBranchName
}
