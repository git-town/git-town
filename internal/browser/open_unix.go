//go:build !windows

package browser

import (
	"github.com/git-town/git-town/v22/internal/filesystem"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// defaultBrowserCommand provides the console command to open the default browser on Unix.
func defaultBrowserCommand() Option[string] {
	return filesystem.FirstExistingExecutable([]string{
		"wsl-open",           // for Windows Subsystem for Linux, see https://github.com/git-town/git-town/issues/1344
		"garcon-url-handler", // opens links in the native browser from Crostini on ChromeOS
		"xdg-open",
		"open",
		"cygstart",
		"x-www-browser",
		"firefox",
		"opera",
		"mozilla",
		"netscape",
	})
}
