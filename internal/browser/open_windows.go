//go:build windows

package browser

import . "github.com/git-town/git-town/v22/pkg/prelude"

// defaultBrowserCommand provides the console command to open the default browser on Windows.
func defaultBrowserCommand() Option[string] {
	// NOTE: the "explorer" command cannot handle special characters like "?" and "=".
	//       In particular, "?" can be escaped via "\", but "=" cannot.
	//       So we are using "start" here.
	return Some("start")
}
