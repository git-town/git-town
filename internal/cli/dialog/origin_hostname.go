package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const (
	originHostnameTitle = `Origin hostname`
	OriginHostnameHelp  = `
If you're using SSH identities,
specify the hostname
of your source code repository.

Only update this
if Git Town's auto-detection doesn't work.

`
)

func OriginHostname(oldValue Option[configdomain.HostingOriginHostname], inputs dialogcomponents.TestInput) (Option[configdomain.HostingOriginHostname], dialogdomain.Exit, error) {
	token, exit, err := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		ExistingValue: oldValue.String(),
		Help:          OriginHostnameHelp,
		Prompt:        "Origin hostname override: ",
		TestInput:     inputs,
		Title:         originHostnameTitle,
	})
	fmt.Printf(messages.OriginHostname, dialogcomponents.FormattedSelection(token, exit))
	return configdomain.ParseHostingOriginHostname(token), exit, err
}
