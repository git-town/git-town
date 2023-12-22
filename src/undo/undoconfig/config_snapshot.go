package undoconfig

import "github.com/git-town/git-town/v11/src/config/configdomain"

// ConfigSnapshot is a snapshot of the Git configuration at a particular point in time.
type ConfigSnapshot struct {
	Global configdomain.SingleCache
	Local  configdomain.SingleCache
}
