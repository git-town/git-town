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
	perennialRegexTitle = `Perennial branch Regex`
	PerennialRegexHelp  = `
Any branch name matching this regular expression
will be treated as a perennial branch.

Example: ^release-.+

If you're not sure what to enter here,
it's safe to leave it blank.

`
	perennialRegexAddendumGlobalOnly = `
The input is prepopulated with the setting
from the global Git configuration.
If you leave it unchanged, Git Town will continue
to use the global setting for this repo.

`
	perennialRegexAddendumGlobalAndLocal = `
Different settings exist in global and local Git metadata.
The input is prepopulated with the local setting.
If you change it to the global one,
Git Town will use the global setting for this repository.

`
)

func PerennialRegex(args PerennialRegexArgs) (Option[configdomain.PerennialRegex], dialogdomain.Exit, error) {
	// if local set --> prepopulate that one, save what the user entered
	// if unscoped set and different from local --> prepopulate that one, save only if different from unscoped
	unscoped, hasUnscoped := args.UnscopedValue.Get()
	local, hasLocal := args.LocalValue.Get()
	hadOnlyLocal := hasLocal && hasUnscoped && local == unscoped
	hadOnlyGlobal := !hasLocal && hasUnscoped
	hadLocalAndGlobal := hasLocal && hasUnscoped && local != unscoped

	helpText := PerennialRegexHelp
	switch {
	case hadOnlyGlobal:
		helpText += perennialRegexAddendumGlobalOnly[1:]
	case hadLocalAndGlobal:
		helpText += perennialRegexAddendumGlobalAndLocal[1:]
	}

	userInputText, exit, err := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		ExistingValue: args.UnscopedValue.String(),
		Help:          helpText,
		Prompt:        "Perennial regex: ",
		TestInput:     args.Input,
		Title:         perennialRegexTitle,
	})
	if err != nil {
		return None[configdomain.PerennialRegex](), false, err
	}
	userInput, err := configdomain.ParsePerennialRegex(userInputText)
	if err != nil {
		return None[configdomain.PerennialRegex](), false, err
	}

	result := None[configdomain.PerennialRegex]()
	switch {
	case hadOnlyLocal:
		result = userInput
	case hadOnlyGlobal && userInput.EqualSome(unscoped):
		result = None[configdomain.PerennialRegex]()
	case hadOnlyGlobal: // user entered a different value than unscoped here
		result = userInput
	case hadLocalAndGlobal && userInput.EqualSome(unscoped): // user entered the global value here
		result = None[configdomain.PerennialRegex]()
	case !hasLocal && !hasUnscoped:
		result = userInput
	}

	fmt.Printf(messages.PerennialRegex, dialogcomponents.FormattedSelection(result.String(), exit))
	return result, exit, nil
}

type PerennialRegexArgs struct {
	Input         dialogcomponents.TestInput
	LocalValue    Option[configdomain.PerennialRegex]
	UnscopedValue Option[configdomain.PerennialRegex]
}
