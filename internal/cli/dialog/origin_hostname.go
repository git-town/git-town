package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
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

func OriginHostname(args Args[configdomain.HostingOriginHostname]) (Option[configdomain.HostingOriginHostname], dialogdomain.Exit, error) {
	input, exit, err := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		DialogName:    "origin-hostname",
		ExistingValue: args.Local.Or(args.Global).StringOr(""),
		Help:          OriginHostnameHelp,
		Inputs:        args.Inputs,
		Prompt:        messages.OriginHostnamePrompt,
		Title:         originHostnameTitle,
	})
	newValue := configdomain.ParseHostingOriginHostname(input)
	if args.Global.Equal(newValue) {
		// the user has entered the global value --> keep using the global value, don't store the local value
		newValue = None[configdomain.HostingOriginHostname]()
	}
	fmt.Printf(messages.OriginHostnameResult, dialogcomponents.FormattedOption(newValue, args.Global.IsSome(), exit))
	return newValue, exit, err
}
