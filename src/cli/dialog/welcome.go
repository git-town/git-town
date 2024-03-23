package dialog

import (
	"github.com/git-town/git-town/v13/src/cli/dialog/components"
)

const (
	welcomeTitle = `Git Town Setup Assistant`
	welcomeText  = `
Welcome to the Git Town setup assistant!
It helps you understand the configuration options for Git Town
and adjust them to match your preferences.

In the following screens, you can change the selection with
UP and DOWN or by entering the entry number.  ENTER goes
to the next screen. Vim motion commands like J, K, O, Q also work.

This assistant only writes changes to disk at the end. You can
try it out safely and exit any time by pressing Q, ESC, or Ctrl-C.

Please press ENTER or O to go to the next screen.

`
)

// MainBranch lets the user select a new main branch for this repo.
func Welcome(inputs components.TestInput) (bool, error) {
	return components.TextDisplay(welcomeTitle, welcomeText, inputs)
}
