package dialog

import (
	"cmp"
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

const (
	branchPrefixTitle = `Branch prefix`
	BranchPrefixHelp  = `
When set, Git Town will automatically
add this prefix to branches it creates or renames.

For example, if set to "feature-",
running "git town hack example"
will create "feature-example".

If you're not sure what to enter here,
it's safe to leave it blank.

`
)

func BranchPrefix(args Args[configdomain.BranchPrefix]) (Option[configdomain.BranchPrefix], dialogdomain.Exit, error) {
	input, exit, errInput := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		DialogName:    "branch-prefix",
		ExistingValue: args.Local.Or(args.Global).StringOr(""),
		Help:          BranchPrefixHelp,
		Inputs:        args.Inputs,
		Prompt:        messages.BranchPrefixPrompt,
		Title:         branchPrefixTitle,
	})
	newValue, errNewValue := configdomain.ParseBranchPrefix(input, "dialog")
	if args.Global.Equal(newValue) {
		// the user has entered the global value --> keep using the global value, don't store the local value
		newValue = None[configdomain.BranchPrefix]()
	}
	fmt.Printf(messages.BranchPrefixResult, dialogcomponents.FormattedOption(newValue, args.Global.IsSome(), exit))
	return newValue, exit, cmp.Or(errInput, errNewValue)
}
