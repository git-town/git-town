package git

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
)

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
