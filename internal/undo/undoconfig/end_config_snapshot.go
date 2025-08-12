package undoconfig

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
)

// ConfigSnapshot is a snapshot of the entire Git configuration at a particular point in time.
type EndConfigSnapshot struct {
	Global configdomain.SingleSnapshot
	Local  configdomain.SingleSnapshot
}
