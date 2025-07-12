package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/colors"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const (
	addendumGlobalOnly = `
The input is prepopulated with the setting
from the global Git configuration.
If you leave it unchanged, Git Town will continue
to use the global setting for this repo.

`
	addendumGlobalAndLocal = `
Different settings exist in global and local Git metadata.
The input is prepopulated with the local setting.
If you change it to the global one,
Git Town will use the global setting for this repository.

`
)

// a part of the setup assistant that allows the user to enter a string-based configuration value
func ConfigStringDialog[T comparable](args ConfigStringDialogArgs[T]) (Option[T], dialogdomain.Exit, error) {
	if args.ConfigFileValue.IsSome() {
		return None[T](), false, nil
	}

	// if local set --> prepopulate that one, save what the user entered
	// if unscoped set and different from local --> prepopulate that one, save only if different from unscoped
	unscoped, hasUnscoped := args.UnscopedValue.Get()
	local, hasLocal := args.LocalValue.Get()
	hadOnlyLocal := hasLocal && hasUnscoped && local == unscoped
	hadOnlyGlobal := !hasLocal && hasUnscoped
	hadLocalAndGlobal := hasLocal && hasUnscoped && local != unscoped

	helpText := args.HelpText
	switch {
	case hadOnlyGlobal:
		helpText += addendumGlobalOnly[1:]
	case hadLocalAndGlobal:
		helpText += addendumGlobalAndLocal[1:]
	}

	var userInput Option[T]
	var exit dialogdomain.Exit
	var err error
	var parseError string
	for {
		var userInputText string
		userInputText, exit, err = dialogcomponents.TextField(dialogcomponents.TextFieldArgs{
			ExistingValue: args.UnscopedValue.String(),
			Help:          helpText + parseError,
			Prompt:        args.Prompt,
			TestInput:     args.Inputs,
			Title:         args.Title,
		})
		if err != nil || exit {
			fmt.Printf(args.ResultMessage, dialogcomponents.FormattedSelection(userInputText, exit))
			return None[T](), exit, err
		}
		userInput, err = args.ParseFunc(userInputText)
		if err != nil {
			parseError = "\n\n" + colors.Red().Styled(err.Error())
			continue
		} else {
			break
		}
	}

	result := None[T]()
	switch {
	case hadOnlyLocal:
		result = userInput
	case hadOnlyGlobal && userInput.EqualSome(unscoped):
		result = None[T]()
	case hadOnlyGlobal: // user entered a different value than unscoped here
		result = userInput
	case hadLocalAndGlobal && userInput.EqualSome(unscoped): // user entered the global value here
		result = None[T]()
	case !hasLocal && !hasUnscoped:
		result = userInput
	}

	fmt.Printf(args.ResultMessage, dialogcomponents.FormattedSelection(result.String(), exit))
	return result, exit, nil
}

type ConfigStringDialogArgs[T any] struct {
	ConfigFileValue Option[T]
	HelpText        string
	Inputs          dialogcomponents.TestInputs
	LocalValue      Option[T]
	ParseFunc       func(string) (Option[T], error)
	Prompt          string
	ResultMessage   string
	Title           string
	UnscopedValue   Option[T]
}
