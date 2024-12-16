package configfile

import (
	"fmt"
	"os"
	"strings"

	"github.com/git-town/git-town/v17/internal/config"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
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
	result.WriteString(fmt.Sprintf("push-new-branches = %t\n", config.NormalConfig.PushNewBranches))
	result.WriteString("\n[hosting]\n")
	result.WriteString(fmt.Sprintf("dev-remote = %q\n", config.NormalConfig.DevRemote.String()))
	if platform, has := config.NormalConfig.HostingPlatform.Get(); has {
		result.WriteString(fmt.Sprintf("platform = %q\n", platform))
	} else {
		result.WriteString("# platform = \"\"\n")
	}
	if config.NormalConfig.HostingOriginHostname.IsNone() {
		result.WriteString("# origin-hostname = \"\"\n")
	} else {
		result.WriteString(fmt.Sprintf("origin-hostname = %q\n", config.NormalConfig.HostingOriginHostname))
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
