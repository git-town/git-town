package bitbucketdatacenter

import "github.com/git-town/git-town/v17/internal/git/giturl"

// Detect always return false because we can't guess a self-hosted URL.
func Detect(_ giturl.Parts) bool {
	return false
}
