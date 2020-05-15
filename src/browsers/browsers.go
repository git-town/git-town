package browsers

import (
	"runtime"

	"github.com/git-town/git-town/src/command"
)

// GetOpenBrowserCommand returns the command to run on the console
// to open the default browser.
func GetOpenBrowserCommand() string {
	if runtime.GOOS == "windows" {
		// NOTE: the "explorer" command cannot handle special characters
		//       like "?" and "=".
		//       In particular, "?" can be escaped via "\", but "=" cannot.
		//       So we are using "start" here.
		return "start"
	}
	for _, browserCommand := range openBrowserCommands {
		res, err := command.Run("which", browserCommand)
		if err == nil && res.OutputSanitized() != "" {
			return browserCommand
		}
	}
	return ""
}

var openBrowserCommands = []string{
	"wsl-open", // for Windows Subsystem for Linux, see https://github.com/git-town/git-town/issues/1344
	"xdg-open",
	"open",
	"cygstart",
	"x-www-browser",
	"firefox",
	"opera",
	"mozilla",
	"netscape",
}
