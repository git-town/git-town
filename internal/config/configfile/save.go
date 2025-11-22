package configfile

import (
	"fmt"
	"os"
	"strings"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
)

func RenderPerennialBranches(perennials gitdomain.LocalBranchNames) string {
	if len(perennials) == 0 {
		return "[]"
	}
	return fmt.Sprintf(`["%s"]`, perennials.Join(`", "`))
}

func RenderTOML(data configdomain.PartialConfig) string {
	result := strings.Builder{}
	result.WriteString("# See https://www.git-town.com/configuration-file for details\n")

	main, hasMain := data.MainBranch.Get()
	displayTypes, hasDisplayTypes := data.DisplayTypes.Get()
	hasPerennialBranches := len(data.PerennialBranches) > 0
	order, hasOrder := data.Order.Get()
	perennialRegex, hasPerennialRegex := data.PerennialRegex.Get()
	if hasMain || hasPerennialBranches || hasPerennialRegex {
		result.WriteString("\n[branches]\n")
		if hasMain {
			result.WriteString(fmt.Sprintf("main = %q\n", main))
		}
		if hasOrder {
			result.WriteString(fmt.Sprintf("order = %q\n", order))
		}
		if hasPerennialBranches {
			result.WriteString(fmt.Sprintf("perennials = %s\n", RenderPerennialBranches(data.PerennialBranches)))
		}
		if hasPerennialRegex {
			result.WriteString(fmt.Sprintf("perennial-regex = %q\n", perennialRegex))
		}
		if hasDisplayTypes {
			result.WriteString(fmt.Sprintf("display-types = %q\n", displayTypes.Serialize(" ")))
		}
	}

	branchPrefix, hasBranchPrefix := data.BranchPrefix.Get()
	newBranchType, hasNewBranchType := data.NewBranchType.Get()
	shareNewBranches, hasShareNewBranches := data.ShareNewBranches.Get()
	stash, hasStash := data.Stash.Get()
	if hasBranchPrefix || hasNewBranchType || hasShareNewBranches || hasStash {
		result.WriteString("\n[create]\n")
		if hasBranchPrefix {
			result.WriteString(fmt.Sprintf("branch-prefix = %q\n", branchPrefix))
		}
		if hasNewBranchType {
			result.WriteString(fmt.Sprintf("new-branch-type = %q\n", newBranchType))
		}
		if hasShareNewBranches {
			result.WriteString(fmt.Sprintf("share-new-branches = %q\n", shareNewBranches))
		}
		if hasStash {
			result.WriteString(fmt.Sprintf("stash = %s\n", stash))
		}
	}

	devRemote, hasDevRemote := data.DevRemote.Get()
	forgeType, hasForgeType := data.ForgeType.Get()
	originHostName, hasOriginHostName := data.HostingOriginHostname.Get()
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

	proposalShowLineage, hasProposalShowLineage := data.ProposalsShowLineage.Get()
	if hasProposalShowLineage {
		result.WriteString("\n[propose]\n")
		result.WriteString(fmt.Sprintf("lineage = %q\n", proposalShowLineage))
	}

	deleteTrackingBranch, hasDeleteTrackingBranch := data.ShipDeleteTrackingBranch.Get()
	shipStrategy, hasShipStrategy := data.ShipStrategy.Get()
	if hasDeleteTrackingBranch || hasShipStrategy {
		result.WriteString("\n[ship]\n")
		if hasDeleteTrackingBranch {
			result.WriteString(fmt.Sprintf("delete-tracking-branch = %t\n", deleteTrackingBranch))
		}
		if hasShipStrategy {
			result.WriteString(fmt.Sprintf("strategy = %q\n", shipStrategy))
		}
	}

	autoResolve, hasAutoResolve := data.AutoResolve.Get()
	autoSync, hasAutoSync := data.AutoSync.Get()
	syncFeatureStrategy, hasSyncFeatureStrategy := data.SyncFeatureStrategy.Get()
	syncPerennialStrategy, hasSyncPerennialStrategy := data.SyncPerennialStrategy.Get()
	syncPrototypeStrategy, hasSyncPrototypeStrategy := data.SyncPrototypeStrategy.Get()
	pushBranches, hasPushBranches := data.PushBranches.Get()
	pushHook, hasPushHook := data.PushHook.Get()
	syncTags, hasSyncTags := data.SyncTags.Get()
	syncUpstream, hasSyncUpstream := data.SyncUpstream.Get()
	if hasAutoResolve || hasSyncFeatureStrategy || hasSyncPerennialStrategy || hasSyncPrototypeStrategy || hasPushHook || hasSyncTags || hasSyncUpstream {
		result.WriteString("\n[sync]\n")
		if hasAutoResolve {
			result.WriteString(fmt.Sprintf("auto-resolve = %t\n", autoResolve))
		}
		if hasAutoSync {
			result.WriteString(fmt.Sprintf("auto-sync = %t\n", autoSync))
		}
		if hasSyncFeatureStrategy {
			result.WriteString(fmt.Sprintf("feature-strategy = %q\n", syncFeatureStrategy))
		}
		if hasSyncPerennialStrategy {
			result.WriteString(fmt.Sprintf("perennial-strategy = %q\n", syncPerennialStrategy))
		}
		if hasSyncPrototypeStrategy {
			result.WriteString(fmt.Sprintf("prototype-strategy = %q\n", syncPrototypeStrategy))
		}
		if hasPushBranches {
			result.WriteString(fmt.Sprintf("push-branches = %t\n", pushBranches))
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

func Save(data configdomain.PartialConfig) error {
	return os.WriteFile(FileName, []byte(RenderTOML(data)), 0o600)
}
