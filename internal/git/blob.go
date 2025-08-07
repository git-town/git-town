package git

import "github.com/git-town/git-town/v21/internal/git/gitdomain"

// describes the content of a file blob in Git
type Blob struct {
	FilePath   string        // relative path of the file in the repo
	Permission string        // permissions, in the form "100755"
	SHA        gitdomain.SHA // checksum of the content blob of the file - this is not the commit SHA!
}

func EmptyBlob() Blob {
	var result Blob
	return result
}
