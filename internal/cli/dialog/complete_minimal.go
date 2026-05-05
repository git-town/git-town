package dialog

import (
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
)

const (
	completeMinimalTitle = `Complete Minimal Setup`
	completeMinimalHelp  = `
Complete the minimal setup of Git Town,
and could end the assistant here
or continue with cli and forge configuration.

You can re-run this assistant: git town init
`
)

func CompleteMinimal(inputs dialogcomponents.Inputs, interactive configdomain.Interactive) (bool, dialogdomain.Exit, error) {
	entries := list.Entries[bool]{
		{
			Data: true,
			Text: `exit and save`,
		},
		{
			Data: false,
			Text: `continue to the cli and forge configuration`,
		},
	}
	selection, exit, err := dialogcomponents.RadioList(entries, 0, completeMinimalTitle, completeMinimalHelp, inputs, interactive, "complete-minimal")
	return selection, exit, err
}
