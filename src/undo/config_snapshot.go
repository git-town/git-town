package undo

import "github.com/git-town/git-town/v11/src/config/gitconfig"

// ConfigSnapshot is a snapshot of the Git configuration at a particular point in time.
type ConfigSnapshot struct {
	GitConfig gitconfig.FullCache
}
