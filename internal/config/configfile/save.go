package configfile

import (
	"fmt"
	"os"
	"strings"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
)

func RenderPerennialBranches(perennials gitdomain.LocalBranchNames) string {
	if len(perennials) == 0 {
		return "[]"
	}
	return fmt.Sprintf(`["%s"]`, perennials.Join(`", "`))
}

func RenderTOML(normalConfig configdomain.NormalConfigData, validatedConfig configdomain.ValidatedConfigData) string {
	result := strings.Builder{}
	result.WriteString("# More info around this file at https://www.git-town.com/configuration-file\n")
	result.WriteString("\n[branches]\n")
	result.WriteString(fmt.Sprintf("main = %q\n", validatedConfig.MainBranch))
	result.WriteString(fmt.Sprintf("perennials = %s\n", RenderPerennialBranches(normalConfig.PerennialBranches)))
	result.WriteString(fmt.Sprintf("perennial-regex = %q\n", normalConfig.PerennialRegex))
	result.WriteString("\n[create]\n")
	result.WriteString(fmt.Sprintf("new-branch-type = %q\n", normalConfig.NewBranchType))
	result.WriteString(fmt.Sprintf("share-new-branches = %q\n", normalConfig.ShareNewBranches))
	result.WriteString("\n[hosting]\n")
	result.WriteString(fmt.Sprintf("dev-remote = %q\n", normalConfig.DevRemote.String()))
	if forgeType, has := normalConfig.ForgeType.Get(); has {
		result.WriteString(fmt.Sprintf("forge-type = %q\n", forgeType))
	}
	if hostName, has := normalConfig.HostingOriginHostname.Get(); has {
		result.WriteString(fmt.Sprintf("origin-hostname = %q\n", hostName))
	}
	result.WriteString("\n[ship]\n")
	result.WriteString(fmt.Sprintf("delete-tracking-branch = %t\n", normalConfig.ShipDeleteTrackingBranch))
	result.WriteString(fmt.Sprintf("strategy = %q\n", normalConfig.ShipStrategy))
	result.WriteString("\n[sync]\n")
	result.WriteString(fmt.Sprintf("feature-strategy = %q\n", normalConfig.SyncFeatureStrategy))
	result.WriteString(fmt.Sprintf("perennial-strategy = %q\n", normalConfig.SyncPerennialStrategy))
	result.WriteString(fmt.Sprintf("prototype-strategy = %q\n", normalConfig.SyncPrototypeStrategy))
	result.WriteString(fmt.Sprintf("push-hook = %t\n", normalConfig.PushHook))
	result.WriteString(fmt.Sprintf("tags = %t\n", normalConfig.SyncTags))
	result.WriteString(fmt.Sprintf("upstream = %t\n", normalConfig.SyncUpstream))
	return result.String()
}

func Save(normalConfig configdomain.NormalConfigData, validatedConfig configdomain.ValidatedConfigData) error {
	return os.WriteFile(FileName, []byte(RenderTOML(normalConfig, validatedConfig)), 0o600)
}
