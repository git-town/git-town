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
	result.WriteString(TOMLComment(dialog.PushHookHelp))
	result.WriteString(fmt.Sprintf("push-hook = %t\n", config.PushHook))
	result.WriteString(TOMLComment(dialog.PushNewBranchesHelp))
	result.WriteString(fmt.Sprintf("push-new-branches = %t\n", config.NewBranchPush))
	result.WriteString(TOMLComment(dialog.ShipDeleteTrackingBranchHelp))
	result.WriteString(fmt.Sprintf("ship-delete-tracking-branch = %t\n", config.ShipDeleteTrackingBranch))
	result.WriteString(TOMLComment(dialog.SyncBeforeShipHelp))
	result.WriteString(fmt.Sprintf("sync-before-ship = %t\n", config.SyncBeforeShip))
	result.WriteString(TOMLComment(dialog.SyncUpstreamHelp))
	result.WriteString(fmt.Sprintf("sync-upstream = %t\n", config.SyncUpstream))
	result.WriteString("\n[branches]\n")
	result.WriteString(TOMLComment(TOMLComment(dialog.MainBranchHelp)))
	result.WriteString(fmt.Sprintf("  main = %s\n", config.MainBranch))
	result.WriteString(TOMLComment(TOMLComment(dialog.PerennialBranchesHelp)))
	result.WriteString(fmt.Sprintf("  perennials = [%s]\n", strings.Join(config.PerennialBranches.Strings(), ", ")))
	result.WriteString("\n[sync-strategy]\n")
	result.WriteString(TOMLComment(TOMLComment(dialog.SyncFeatureStrategyHelp)))
	result.WriteString(fmt.Sprintf("  feature-branches = %s\n", config.SyncFeatureStrategy))
	result.WriteString(TOMLComment(TOMLComment(dialog.SyncPerennialStrategyHelp)))
	result.WriteString(fmt.Sprintf("  perennial-branches = %s\n", config.SyncPerennialStrategy))
	return result.String()
}

func TOMLComment(text string) string {
	if text == "" {
		return ""
	}
	result := []string{}
	for _, line := range strings.Split(text, "\n") {
		result = append(result, "# "+line)
	}
	return strings.Join(result, "\n")
}
