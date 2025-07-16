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

func RenderTOML(config configdomain.PartialConfig) string {
	result := strings.Builder{}
	result.WriteString("# More info around this file at https://www.git-town.com/configuration-file\n")

	main, hasMain := config.MainBranch.Get()
	hasPerennialBranches := len(config.PerennialBranches) > 0
	perennialRegex, hasPerennialRegex := config.PerennialRegex.Get()
	if hasMain || hasPerennialBranches || hasPerennialRegex {
		result.WriteString("\n[branches]\n")
		if hasMain {
			result.WriteString(fmt.Sprintf("main = %q\n", main))
		}
		if hasPerennialBranches {
			result.WriteString(fmt.Sprintf("perennials = %s\n", RenderPerennialBranches(config.PerennialBranches)))
		}
		if hasPerennialRegex {
			result.WriteString(fmt.Sprintf("perennial-regex = %q\n", perennialRegex))
		}
	}

	newBranchType, hasNewBranchType := config.NewBranchType.Get()
	shareNewBranches, hasShareNewBranches := config.ShareNewBranches.Get()
	if hasNewBranchType || hasShareNewBranches {
		result.WriteString("\n[create]\n")
		if hasNewBranchType {
			result.WriteString(fmt.Sprintf("new-branch-type = %q\n", newBranchType))
		}
		if hasShareNewBranches {
			result.WriteString(fmt.Sprintf("share-new-branches = %q\n", shareNewBranches))
		}
	}

	devRemote, hasDevRemote := config.DevRemote.Get()
	forgeType, hasForgeType := config.ForgeType.Get()
	originHostName, hasOriginHostName := config.HostingOriginHostname.Get()
	if hasDevRemote || hasForgeType || hasOriginHostName {
		result.WriteString("\n[hosting]\n")
		if hasDevRemote {
			result.WriteString(fmt.Sprintf("dev-remote = %q\n", devRemote))
		}
		if hasForgeType {
			result.WriteString(fmt.Sprintf("forge-type = %q\n", forgeType))
		}
		if hasOriginHostName {
			result.WriteString(fmt.Sprintf("origin-hostname = %q\n", originHostName))
		}
	}

	deleteTrackingBranch, hasDeleteTrackingBranch := config.ShipDeleteTrackingBranch.Get()
	shipStrategy, hasShipStrategy := config.ShipStrategy.Get()
	if hasDeleteTrackingBranch || hasShipStrategy {
		result.WriteString("\n[ship]\n")
		if hasDeleteTrackingBranch {
			result.WriteString(fmt.Sprintf("delete-tracking-branch = %t\n", deleteTrackingBranch))
		}
		if hasShipStrategy {
			result.WriteString(fmt.Sprintf("strategy = %q\n", shipStrategy))
		}
	}

	syncFeatureStrategy, hasSyncFeatureStrategy := config.SyncFeatureStrategy.Get()
	syncPerennialStrategy, hasSyncPerennialStrategy := config.SyncPerennialStrategy.Get()
	syncPrototypeStrategy, hasSyncPrototypeStrategy := config.SyncPrototypeStrategy.Get()
	pushHook, hasPushHook := config.PushHook.Get()
	syncTags, hasSyncTags := config.SyncTags.Get()
	syncUpstream, hasSyncUpstream := config.SyncUpstream.Get()
	if hasSyncFeatureStrategy || hasSyncPerennialStrategy || hasSyncPrototypeStrategy || hasPushHook || hasSyncTags || hasSyncUpstream {
		result.WriteString("\n[sync]\n")
		if hasSyncFeatureStrategy {
			result.WriteString(fmt.Sprintf("feature-strategy = %q\n", syncFeatureStrategy))
		}
		if hasSyncPerennialStrategy {
			result.WriteString(fmt.Sprintf("perennial-strategy = %q\n", syncPerennialStrategy))
		}
		if hasSyncPrototypeStrategy {
			result.WriteString(fmt.Sprintf("prototype-strategy = %q\n", syncPrototypeStrategy))
		}
		if hasPushHook {
			result.WriteString(fmt.Sprintf("push-hook = %t\n", pushHook))
		}
		if hasSyncTags {
			result.WriteString(fmt.Sprintf("tags = %t\n", syncTags))
		}
		if hasSyncUpstream {
			result.WriteString(fmt.Sprintf("upstream = %t\n", syncUpstream))
		}
	}
	return result.String()
}

func Save(config configdomain.PartialConfig) error {
	return os.WriteFile(FileName, []byte(RenderTOML(config)), 0o600)
}
