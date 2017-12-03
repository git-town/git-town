package browsers

import (
	"runtime"

	"github.com/Originate/git-town/src/command"
	"github.com/Originate/git-town/src/util"
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
		cmd := command.New("which", browserCommand)
		if cmd.Err() == nil && cmd.Output() != "" {
			return browserCommand
		}
	}
	util.ExitWithErrorMessage(missingOpenBrowserCommandMessages...)
	return ""
}

var openBrowserCommands = []string{
	"xdg-open",
	"open",
	"cygstart",
	"x-www-browser",
	"firefox",
	"opera",
	"mozilla",
	"netscape",
}

var missingOpenBrowserCommandMessages = []string{
	"Cannot open a browser.",
	"If you think this is a bug,",
	"please open an issue at https://github.com/Originate/git-town/issues",
	"and mention your OS and browser.",
}
