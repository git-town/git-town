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
)

func PerennialRegex(args PerennialRegexArgs) (Option[configdomain.PerennialRegex], dialogdomain.Exit, error) {
	// if local set --> prepopulate that one, save what the user entered
	// if unscoped set and different from local --> prepopulate that one, save only if different from unscoped
	unscoped, hasUnscoped := args.UnscopedValue.Get()
	local, hasLocal := args.LocalValue.Get()
	hadOnlyLocal := hasLocal && hasUnscoped && local == unscoped
	hadOnlyGlobal := !hasLocal && hasUnscoped
	hadLocalAndGlobal := hasLocal && hasUnscoped && local != unscoped

	// TODO: add instructions to user to PerennialRegexHelp

	userInputText, exit, err := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		ExistingValue: args.UnscopedValue.String(),
		Help:          PerennialRegexHelp,
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
