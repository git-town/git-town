package update

import (
	"context"
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/pkg/colors"
)

// Notifier handles displaying update notifications to users
type Notifier struct {
	logger         print.Logger
	versionChecker *VersionChecker
}

// NewNotifier creates a new update notifier
func NewNotifier(logger print.Logger) *Notifier {
	return &Notifier{
		logger:         logger,
		versionChecker: NewVersionChecker(),
	}
}

// CheckAndNotify checks for updates and displays a notification if available
func (self *Notifier) CheckAndNotify(ctx context.Context) {
	updateInfo, err := self.versionChecker.CheckForUpdate(ctx)
	if err != nil {
		return
	}

	if updateInfo.IsUpdateAvailable() {
		self.displayUpdateNotification(updateInfo)
	}
}

// displayUpdateNotification shows the update notification to the user
func (self *Notifier) displayUpdateNotification(updateInfo *Info) {
	bold := colors.Bold()
	cyan := colors.BoldCyan()
	green := colors.BoldGreen()

	fmt.Println()
	fmt.Printf("%s\n", bold.Styled("╭─────────────────────────────────────────────────────────────╮"))
	fmt.Printf("%s\n", bold.Styled("│                    🚀 Update Available!                    │"))
	fmt.Printf("%s\n", bold.Styled("├─────────────────────────────────────────────────────────────┤"))
	fmt.Printf("%s %s → %s\n",
		bold.Styled("│ Git Town:"),
		cyan.Styled("v"+updateInfo.CurrentVersion),
		green.Styled("v"+updateInfo.LatestVersion))
	fmt.Printf("%s\n", bold.Styled("│                                                             │"))
	fmt.Printf("%s %s\n",
		bold.Styled("│ Download:"),
		cyan.Styled(updateInfo.UpdateURL))
	fmt.Printf("%s\n", bold.Styled("╰─────────────────────────────────────────────────────────────╯"))
	fmt.Println()
}
