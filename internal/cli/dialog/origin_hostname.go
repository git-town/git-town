package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v19/internal/cli/dialog/components"
	"github.com/git-town/git-town/v19/internal/config/configdomain"
	"github.com/git-town/git-town/v19/internal/messages"
	. "github.com/git-town/git-town/v19/pkg/prelude"
)

const (
	originHostnameTitle = `Origin hostname`
	OriginHostnameHelp  = `
If you're using SSH identities,
specify the hostname of your source code repository.

Only update this if Git Town's auto-detection doesn't work.

`
)

func OriginHostname(oldValue Option[configdomain.HostingOriginHostname], inputs components.TestInput) (Option[configdomain.HostingOriginHostname], bool, error) {
	token, aborted, err := components.TextField(components.TextFieldArgs{
		ExistingValue: oldValue.String(),
		Help:          OriginHostnameHelp,
		Prompt:        "Origin hostname override: ",
		TestInput:     inputs,
		Title:         originHostnameTitle,
	})
	fmt.Printf(messages.OriginHostname, components.FormattedToken(token, aborted))
	return configdomain.ParseHostingOriginHostname(token), aborted, err
}
