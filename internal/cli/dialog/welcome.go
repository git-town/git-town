package dialog

import (
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
)

const (
	welcomeTitle = `Git Town Setup Assistant`
	welcomeText  = `
Welcome to the Git Town setup assistant!
This tool walks you through the available configuration options for Git Town
and helps you tailor them to your workflow.

On the next screens, navigate using the UP/DOWN arrows
or by typing the corresponding entry number.
Press ENTER to proceed.
Vim-style motions like J, K, O, and Q are also supported.

No changes are written until the final step,
so feel free to explore.
You can exit at any time with Q, ESC, or Ctrl-C.

You can execute this setup assistant again by running
"git town init" in your terminal.

When you're ready, press ENTER or O to continue.

`
)

func Welcome(inputs dialogcomponents.Inputs) (dialogdomain.Exit, error) {
	return dialogcomponents.TextDisplay(welcomeTitle, welcomeText, inputs, "welcome")
}
