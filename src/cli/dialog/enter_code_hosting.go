package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const enterHostingPlatformHelp = `
Git Town needs to know on which code hosting platform you host your code.
This allows it to open the right browser URLs and talk to the right API endpoints.
Only change this setting if the auto-detection does not work for you.

`

// EnterMainBranch lets the user select a new main branch for this repo.
func EnterHostingPlatform(platformName configdomain.CodeHostingPlatformName, inputs TestInput) (configdomain.CodeHostingPlatformName, bool, error) {
	selection, aborted, err := radioList(radioListArgs{
		entries: []string{
			configdomain.CodeHostingPlatformNameAutoDetect,
			configdomain.CodeHostingPlatformBitBucket,
			configdomain.CodeHostingPlatformGitea,
			configdomain.CodeHostingPlatformGitHub,
			configdomain.CodeHostingPlatformGitLab,
		},
		defaultEntry: platformName.String(),
		help:         enterHostingPlatformHelp,
		testInput:    inputs,
	})
	fmt.Printf("Code hosting: %s\n", formattedSelection(selection, aborted))
	return configdomain.NewCodeHostingPlatformName(selection), aborted, err
}
