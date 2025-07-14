package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/colors"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type comparableStringer interface {
	comparable
	fmt.Stringer
}

// a part of the setup assistant that allows the user to enter a string-based configuration value
func ConfigEnumDialog[T comparableStringer](args ConfigEnumDialogArgs[T]) (Option[T], dialogdomain.Exit, error) {
	if args.ConfigFileValue.IsSome() {
		return None[T](), false, nil
	}

	// local and same as unscoped --> select that one, save what the user entered
	// !local and uncoped --> prepopulate the unscoped
	// if local and unscoped
	// if unscoped set and different from local --> prepopulate that one, save only if different from unscoped
	unscoped, hasUnscoped := args.UnscopedValue.Get()
	local, hasLocal := args.LocalValue.Get()
	hadOnlyLocal := hasLocal && hasUnscoped && local == unscoped
	hadOnlyGlobal := !hasLocal && hasUnscoped
	hadLocalAndGlobal := hasLocal && hasUnscoped && local != unscoped

	defaultPos := 0
	helpText := args.HelpText
	entries := args.Entries
	switch {
	case hadOnlyGlobal:
		entries = appendSkipEntry(entries)
		defaultPos = 1
		helpText += addendumGlobalOnly[1:]
	case hadLocalAndGlobal:
		helpText += addendumGlobalAndLocal[1:]
	}

	var userInput Option[T]
	var exit dialogdomain.Exit
	var parseError string
	for {
		selection, exit, err := dialogcomponents.RadioList(args.Entries, defaultPos, args.Title, helpText+parseError, args.Inputs.Next())
		fmt.Printf(messages.SyncFeatureBranches, dialogcomponents.FormattedSelection(selection.String(), exit))
		// return selection, exit, err
		if err != nil || exit {
			return None[T](), exit, err
		}
		userInput, err = args.ParseFunc(selection.String())
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

type ConfigEnumDialogArgs[T comparable] struct {
	ConfigFileValue Option[T]
	Entries         list.Entries[T]
	HelpText        string
	Inputs          dialogcomponents.TestInputs
	LocalValue      Option[T]
	ParseFunc       func(string) (Option[T], error)
	Prompt          string
	ResultMessage   string
	Title           string
	UnscopedValue   Option[T]
}
