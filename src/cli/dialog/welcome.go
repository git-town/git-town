package dialog

import (
	"github.com/git-town/git-town/v11/src/cli/dialog/components"
)

const welcomeText = `
Welcome to the Git Town setup assistant! It will walk you
step by step through the configuration options for Git Town.
It explains what each option does and lets you modify it.

You can use the UP and DOWN keys to change the selection and
ENTER to accept it and go to the next screen.  Vim motion keys
("j" and "k" for DOWN and UP and "o" to accept) also work!

I only persist changes at the end. You can try this assistant
and exit safely at any time by pressing ESC or Ctrl-C.
Please press ENTER or "o" to go to the next screen.

`

// MainBranch lets the user select a new main branch for this repo.
func Welcome(inputs components.TestInput) (bool, error) {
	return components.TextDisplay(welcomeText, inputs)
}
