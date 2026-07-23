package format

import (
	"fmt"

	"github.com/git-town/git-town/v24/internal/config/configdomain"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

// Secret provides a printable version of the given secret configuration value.
// Configured secrets are redacted as "(configured)" unless showSecrets is enabled.
func Secret[T fmt.Stringer](secret Option[T], showSecrets configdomain.ShowSecrets) string {
	if !showSecrets.ShouldShowSecrets() && secret.IsSome() {
		return "(configured)"
	}
	return OptionalStringerSetting(secret)
}
