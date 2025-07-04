package configfile

import (
	"fmt"
	"os"
	"strings"

	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
)

func RenderPerennialBranches(perennials gitdomain.LocalBranchNames) string {
	if len(perennials) == 0 {
		return "[]"
	}
	return fmt.Sprintf(`["%s"]`, perennials.Join(`", "`))
}

func RenderTOML(config *config.UnvalidatedConfig) string {
	result := strings.Builder{}
	result.WriteString("# More info around this file at https://www.git-town.com/configuration-file\n")
	result.WriteString("\n[branches]\n")
	result.WriteString(fmt.Sprintf("main = %q\n", config.UnvalidatedConfig.MainBranch))
	result.WriteString(fmt.Sprintf("perennials = %s\n", RenderPerennialBranches(config.NormalConfig.PerennialBranches)))
	result.WriteString(fmt.Sprintf("perennial-regex = %q\n", config.NormalConfig.PerennialRegex))
	result.WriteString("\n[create]\n")
	result.WriteString(fmt.Sprintf("new-branch-type = %q\n", config.NormalConfig.NewBranchType))
	result.WriteString(fmt.Sprintf("share-new-branches = %q\n", config.NormalConfig.ShareNewBranches))
	result.WriteString("\n[hosting]\n")
	result.WriteString(fmt.Sprintf("dev-remote = %q\n", config.NormalConfig.DevRemote.String()))
	if forgeType, has := config.NormalConfig.ForgeType.Get(); has {
		result.WriteString(fmt.Sprintf("forge-type = %q\n", forgeType))
	}
	if hostName, has := config.NormalConfig.HostingOriginHostname.Get(); has {
		result.WriteString(fmt.Sprintf("origin-hostname = %q\n", hostName))
	}
	result.WriteString("\n[ship]\n")
	result.WriteString(fmt.Sprintf("delete-tracking-branch = %t\n", config.NormalConfig.ShipDeleteTrackingBranch))
	result.WriteString(fmt.Sprintf("strategy = %q\n", config.NormalConfig.ShipStrategy))
	result.WriteString("\n[sync]\n")
	result.WriteString(fmt.Sprintf("feature-strategy = %q\n", config.NormalConfig.SyncFeatureStrategy))
	result.WriteString(fmt.Sprintf("perennial-strategy = %q\n", config.NormalConfig.SyncPerennialStrategy))
	result.WriteString(fmt.Sprintf("prototype-strategy = %q\n", config.NormalConfig.SyncPrototypeStrategy))
	result.WriteString(fmt.Sprintf("push-hook = %t\n", config.NormalConfig.PushHook))
	result.WriteString(fmt.Sprintf("tags = %t\n", config.NormalConfig.SyncTags))
	result.WriteString(fmt.Sprintf("upstream = %t\n", config.NormalConfig.SyncUpstream))
	return result.String()
}

func Save(config *config.UnvalidatedConfig) error {
	return os.WriteFile(FileName, []byte(RenderTOML(config)), 0o600)
}
