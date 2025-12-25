package dialog

import (
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
)

const (
	enterAllTitle = `Go through all configuration options?`
	enterAllHelp  = `
You are good to go with the basic setup of Git Town,
and could end the assistant here
or continue with the advanced configuration options.

You can re-run this assistant: git town init
`
)

func EnterAll(inputs dialogcomponents.Inputs) (bool, dialogdomain.Exit, error) {
	entries := list.Entries[bool]{
		{
			Data: false,
			Text: `exit and save`,
		},
		{
			Data: true,
			Text: `continue to the advanced configuration options`,
		},
	}
	selection, exit, err := dialogcomponents.RadioList(entries, 0, enterAllTitle, enterAllHelp, inputs, "enter-all")
	return selection, exit, err
}
