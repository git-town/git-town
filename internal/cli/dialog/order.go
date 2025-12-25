package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

const (
	orderTitle = `Branch ordering`
	orderHelp  = `
How should Git Town order branches it displays?

`
)

func Order(args Args[configdomain.Order]) (Option[configdomain.Order], dialogdomain.Exit, error) {
	entries := list.Entries[Option[configdomain.Order]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[configdomain.Order]]{
			Data: None[configdomain.Order](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}
	entries = append(entries, list.Entries[Option[configdomain.Order]]{
		{
			Data: Some(configdomain.OrderAsc),
			Text: "natural sorting (ascending)",
		},
		{
			Data: Some(configdomain.OrderDesc),
			Text: "natural sorting (descending)",
		},
	}...)
	defaultPos := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, orderTitle, orderHelp, args.Inputs, "order")
	fmt.Printf(messages.Order, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}
