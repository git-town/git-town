package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
)

const (
	originHostnameTitle = `Origin hostname`
	OriginHostnameHelp  = `
When using SSH identities, define the hostname
of your source code repository. Only change this
if the auto-detection does not work for you.

`
)

// GitHubToken lets the user enter the GitHub API token.
func OriginHostname(oldValue Option[configdomain.HostingOriginHostname], inputs components.TestInput) (Option[configdomain.HostingOriginHostname], bool, error) {
	token, aborted, err := components.TextField(components.TextFieldArgs{
		ExistingValue: oldValue.String(),
		Help:          OriginHostnameHelp,
		Prompt:        "Origin hostname override: ",
		TestInput:     inputs,
		Title:         originHostnameTitle,
	})
	fmt.Printf(messages.OriginHostname, components.FormattedToken(token, aborted))
	return configdomain.NewHostingOriginHostnameOption(token), aborted, err
}
