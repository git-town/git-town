package dialog

import (
	"github.com/git-town/git-town/v11/src/cli/dialog/components"
)

const (
	welcomeTitle = `Git Town Setup Assistant`
	welcomeText  = `
Welcome to the Git Town setup assistant!
It helps you configure Git Town for your system.

In the following screens, you can change the selection with
UP and DOWN or by entering the entry number.  ENTER goes
to the next screen. Vim motion commands like j, k, o also work.

This assistant only writes changes to disk at the end. You can
try it out safely and exit any time by pressing ESC, Ctrl-C, or q.

Please press ENTER or "o" to go to the next screen.

`
)

// MainBranch lets the user select a new main branch for this repo.
func Welcome(inputs components.TestInput) (bool, error) {
	return components.TextDisplay(welcomeTitle, welcomeText, inputs)
}
