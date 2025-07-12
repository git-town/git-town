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
	featureRegexTitle = `Feature branch regex`
	FeatureRegexHelp  = `
Branches matching this regular expression
will be treated as feature branches.
This setting only applies
if the "unknown-branch-type"
is set to something other than "feature".

`
)

func FeatureRegex(args FeatureRegexArgs) (Option[configdomain.FeatureRegex], dialogdomain.Exit, error) {
	// if local set --> prepopulate that one, save what the user entered
	// if unscoped set and different from local --> prepopulate that one, save only if different from unscoped
	unscoped, hasUnscoped := args.UnscopedValue.Get()
	local, hasLocal := args.LocalValue.Get()
	hadOnlyLocal := hasLocal && hasUnscoped && local == unscoped
	hadOnlyGlobal := !hasLocal && hasUnscoped
	hadLocalAndGlobal := hasLocal && hasUnscoped && local != unscoped

	helpText := FeatureRegexHelp
	switch {
	case hadOnlyGlobal:
		helpText += addendumGlobalOnly[1:]
	case hadLocalAndGlobal:
		helpText += addendumGlobalAndLocal[1:]
	}

	userInputText, exit, err := dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
		ExistingValue: args.UnscopedValue.String(),
		Help:          helpText,
		Prompt:        "Feature regex: ",
		TestInput:     args.Input,
		Title:         featureRegexTitle,
	})
	if err != nil {
		return None[configdomain.FeatureRegex](), false, err
	}
	userInput, err := configdomain.ParseFeatureRegex(userInputText)
	if err != nil {
		return None[configdomain.FeatureRegex](), false, err
	}

	result := None[configdomain.FeatureRegex]()
	switch {
	case hadOnlyLocal:
		result = userInput
	case hadOnlyGlobal && userInput.EqualSome(unscoped):
		result = None[configdomain.FeatureRegex]()
	case hadOnlyGlobal: // user entered a different value than unscoped here
		result = userInput
	case hadLocalAndGlobal && userInput.EqualSome(unscoped): // user entered the global value here
		result = None[configdomain.FeatureRegex]()
	case !hasLocal && !hasUnscoped:
		result = userInput
	}

	fmt.Printf(messages.FeatureRegex, dialogcomponents.FormattedSelection(result.String(), exit))
	return result, exit, nil
}

type FeatureRegexArgs struct {
	Input         dialogcomponents.TestInput
	LocalValue    Option[configdomain.FeatureRegex]
	UnscopedValue Option[configdomain.FeatureRegex]
}
