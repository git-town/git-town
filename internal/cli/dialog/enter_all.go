package dialog

import (
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
)

const (
	enterAllTitle = `Go through all configuration options?`
	enterAllHelp  = `
You are good to go with a basic setup of Git Town
and could end the assistant here
or go through all available configuration options and learn about them.

You can always re-run this assistant later through "git town config setup".
`
)

func EnterAll(inputs dialogcomponents.Inputs) (bool, dialogdomain.Exit, error) {
	entries := list.Entries[bool]{
		{
			Data: false,
			Text: `exit here`,
		},
		{
			Data: true,
			Text: `go through all configuration options`,
		},
	}
	selection, exit, err := dialogcomponents.RadioList(entries, 0, enterAllTitle, enterAllHelp, inputs, "enter-all")
	return selection, exit, err
}
