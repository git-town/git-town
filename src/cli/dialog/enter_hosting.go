package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/gohacks/stringers"
)

const enterHostingPlatformHelp = `
Git Town needs to know on which code hosting platform you host your code.
This allows it to open the right browser URLs and talk to the right API endpoints.
Only change this setting if the auto-detection does not work for you.

`

// EnterMainBranch lets the user select a new main branch for this repo.
func EnterHostingPlatform(platformName configdomain.HostingPlatform, inputs TestInput) (configdomain.HostingPlatform, bool, error) {
	entries := []configdomain.HostingPlatform{
		configdomain.HostingPlatformAutoDetect,
		configdomain.HostingPlatformBitBucket,
		configdomain.HostingPlatformGitea,
		configdomain.HostingPlatformGitHub,
		configdomain.HostingPlatformGitLab,
	}
	cursor := stringers.IndexOrStart(entries, platformName)
	selection, aborted, err := radioList(entries, cursor, enterHostingPlatformHelp, inputs)
	fmt.Printf("Code hosting: %s\n", formattedSelection(selection.String(), aborted))
	return selection, aborted, err
}
