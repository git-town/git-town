package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/dialog/components"
	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const originHostnameHelp = `
When using SSH identities, define the hostname of your source code repository.

Only change this setting if the auto-detection does not work for you.

`

// GitHubToken lets the user enter the GitHub API token.
func OriginHostname(oldValue configdomain.HostingOriginHostname, inputs components.TestInput) (configdomain.HostingOriginHostname, bool, error) {
	token, aborted, err := components.TextField(components.TextFieldArgs{
		ExistingValue: oldValue.String(),
		Help:          originHostnameHelp,
		Prompt:        "Origin hostname override: ",
		TestInput:     inputs,
	})
	fmt.Printf("Origin hostname: %s\n", components.FormattedToken(token, aborted))
	return configdomain.HostingOriginHostname(token), aborted, err
}
