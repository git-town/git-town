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

func RenderTOML(data configdomain.PartialConfig, mainBranch gitdomain.LocalBranchName) string {
	result := strings.Builder{}
	result.WriteString("# More info around this file at https://www.git-town.com/configuration-file\n")
	result.WriteString("\n[branches]\n")
	result.WriteString(fmt.Sprintf("main = %q\n", mainBranch))
	if len(data.PerennialBranches) > 0 {
		result.WriteString(fmt.Sprintf("perennials = %s\n", RenderPerennialBranches(data.PerennialBranches)))
	}
	if value, has := data.PerennialRegex.Get(); has {
		result.WriteString(fmt.Sprintf("perennial-regex = %q\n", value))
	}
	result.WriteString("\n[create]\n")
	result.WriteString(fmt.Sprintf("new-branch-type = %q\n", data.NewBranchType))
	result.WriteString(fmt.Sprintf("share-new-branches = %q\n", data.ShareNewBranches))
	result.WriteString("\n[hosting]\n")
	result.WriteString(fmt.Sprintf("dev-remote = %q\n", data.DevRemote.String()))
	if forgeType, has := data.ForgeType.Get(); has {
		result.WriteString(fmt.Sprintf("forge-type = %q\n", forgeType))
	}
	if hostName, has := data.HostingOriginHostname.Get(); has {
		result.WriteString(fmt.Sprintf("origin-hostname = %q\n", hostName))
	}
	result.WriteString("\n[ship]\n")
	result.WriteString(fmt.Sprintf("delete-tracking-branch = %t\n", data.ShipDeleteTrackingBranch))
	result.WriteString(fmt.Sprintf("strategy = %q\n", data.ShipStrategy))
	result.WriteString("\n[sync]\n")
	result.WriteString(fmt.Sprintf("feature-strategy = %q\n", data.SyncFeatureStrategy))
	result.WriteString(fmt.Sprintf("perennial-strategy = %q\n", data.SyncPerennialStrategy))
	result.WriteString(fmt.Sprintf("prototype-strategy = %q\n", data.SyncPrototypeStrategy))
	result.WriteString(fmt.Sprintf("push-hook = %t\n", data.PushHook))
	result.WriteString(fmt.Sprintf("tags = %t\n", data.SyncTags))
	result.WriteString(fmt.Sprintf("upstream = %t\n", data.SyncUpstream))
	return result.String()
}

func Save(normalConfig configdomain.PartialConfig, mainBranch gitdomain.LocalBranchName) error {
	return os.WriteFile(FileName, []byte(RenderTOML(normalConfig, mainBranch)), 0o600)
}
