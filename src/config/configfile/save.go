package configfile

import (
	"fmt"
	"os"
	"strings"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/config/configdomain"
)

func Save(config *configdomain.FullConfig) error {
	toml := RenderTOML(config)
	return os.WriteFile(FileName, []byte(toml), 0o600)
}

func RenderTOML(config *configdomain.FullConfig) string {
	result := strings.Builder{}
	result.WriteString("# Git Town configuration file\n#\n")
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.PushHookHelp), "") + "\n")
	result.WriteString(fmt.Sprintf("push-hook = %t\n\n", config.PushHook))
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.PushNewBranchesHelp), "") + "\n")
	result.WriteString(fmt.Sprintf("push-new-branches = %t\n\n", config.NewBranchPush))
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.ShipDeleteTrackingBranchHelp), "") + "\n")
	result.WriteString(fmt.Sprintf("ship-delete-tracking-branch = %t\n\n", config.ShipDeleteTrackingBranch))
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.SyncBeforeShipHelp), "") + "\n")
	result.WriteString(fmt.Sprintf("sync-before-ship = %t\n\n", config.SyncBeforeShip))
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.SyncUpstreamHelp), "") + "\n")
	result.WriteString(fmt.Sprintf("sync-upstream = %t\n", config.SyncUpstream))
	result.WriteString("\n[branches]\n")
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.MainBranchHelp), "  ") + "\n")
	result.WriteString(fmt.Sprintf("  main = %q\n\n", config.MainBranch))
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.PerennialBranchesHelp), "  ") + "\n")
	result.WriteString(fmt.Sprintf("  perennials = [\"%s\"]\n", config.PerennialBranches.Join(`", "`)))
	result.WriteString("\n[sync-strategy]\n")
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.SyncFeatureStrategyHelp), "  ") + "\n")
	result.WriteString(fmt.Sprintf("  feature-branches = %q\n\n", config.SyncFeatureStrategy))
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.SyncPerennialStrategyHelp), "  ") + "\n")
	result.WriteString(fmt.Sprintf("  perennial-branches = %q\n", config.SyncPerennialStrategy))
	return result.String()
}

func TOMLComment(text string, linePrefix string) string {
	if text == "" {
		return ""
	}
	result := []string{}
	for _, line := range strings.Split(text, "\n") {
		if line == "" {
			result = append(result, linePrefix+"#")
		} else {
			result = append(result, linePrefix+"# "+line)
		}
	}
	return strings.Join(result, "\n")
}
