package git

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
)

// Blob describes the content of a file blob in Git.
type Blob struct {
	FilePath   string        // relative path of the file in the repo
	Permission string        // permissions, in the form "100755"
	SHA        gitdomain.SHA // checksum of the content blob of the file - this is not the commit SHA!
}

func (bi Blob) Debug(querier subshelldomain.Querier) {
	fileContent, err := querier.Query("git", "show", bi.SHA.String())
	if err != nil {
		panic(fmt.Sprintf("cannot display content of blob %q: %s", bi.SHA, err))
	}
	fmt.Printf("%s %s %s\n%s", bi.FilePath, bi.SHA.Truncate(7), bi.Permission, gohacks.IndentLines(fileContent, 4))
}

func EmptyBlob() Blob {
	var result Blob
	return result
}
