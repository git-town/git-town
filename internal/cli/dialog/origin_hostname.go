package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v18/internal/cli/dialog/components"
	"github.com/git-town/git-town/v18/internal/config/configdomain"
	"github.com/git-town/git-town/v18/internal/messages"
	. "github.com/git-town/git-town/v18/pkg/prelude"
)

const (
	originHostnameTitle = `Origin hostname`
	OriginHostnameHelp  = `
When using SSH identities, define the hostname
of your source code repository. Only change this
if the auto-detection does not work for you.

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
