package configfile

import (
	"fmt"
	"os"
	"strings"

	"github.com/git-town/git-town/v15/internal/cli/dialog"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
)

func RenderPerennialBranches(perennials gitdomain.LocalBranchNames) string {
	if len(perennials) == 0 {
		return "[]"
	}
	return fmt.Sprintf(`["%s"]`, perennials.Join(`", "`))
}

func RenderTOML(config *configdomain.UnvalidatedConfig) string {
	result := strings.Builder{}
	result.WriteString("# Git Town configuration file\n")
	result.WriteString("#\n")
	result.WriteString("# Run \"git town config setup\" to add additional entries\n")
	result.WriteString("# to this file after updating Git Town.\n")
	result.WriteString("#\n")
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.PushHookHelp)) + "\n")
	result.WriteString(fmt.Sprintf("push-hook = %t\n\n", config.PushHook))
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.PushNewBranchesHelp)) + "\n")
	result.WriteString(fmt.Sprintf("push-new-branches = %t\n\n", config.PushNewBranches))
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.CreatePrototypeBranchesHelp)) + "\n")
	result.WriteString(fmt.Sprintf("create-prototype-branches = %t\n\n", config.CreatePrototypeBranches))
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.ShipDeleteTrackingBranchHelp)) + "\n")
	result.WriteString(fmt.Sprintf("ship-delete-tracking-branch = %t\n\n", config.ShipDeleteTrackingBranch))
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.SyncTagsHelp)) + "\n")
	result.WriteString(fmt.Sprintf("sync-tags = %t\n\n", config.SyncTags))
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.SyncUpstreamHelp)) + "\n")
	result.WriteString(fmt.Sprintf("sync-upstream = %t\n", config.SyncUpstream))
	result.WriteString("\n[branches]\n\n")
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.MainBranchHelp)) + "\n")
	result.WriteString(fmt.Sprintf("main = %q\n\n", config.MainBranch))
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.PerennialBranchesHelp)) + "\n")
	result.WriteString(fmt.Sprintf("perennials = %s\n", RenderPerennialBranches(config.PerennialBranches)) + "\n")
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.PerennialRegexHelp)) + "\n")
	result.WriteString(fmt.Sprintf("perennial-regex = %q\n", config.PerennialRegex))
	result.WriteString("\n[hosting]\n\n")
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.HostingPlatformHelp)) + "\n")
	if platform, has := config.HostingPlatform.Get(); has {
		result.WriteString(fmt.Sprintf("platform = %q\n\n", platform))
	} else {
		result.WriteString("# platform = \"\"\n\n")
	}
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.OriginHostnameHelp)) + "\n")
	if config.HostingOriginHostname.IsNone() {
		result.WriteString("# origin-hostname = \"\"\n")
	} else {
		result.WriteString(fmt.Sprintf("origin-hostname = %q\n", config.HostingOriginHostname))
	}
	result.WriteString("\n[sync-strategy]\n\n")
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.SyncFeatureStrategyHelp)) + "\n")
	result.WriteString(fmt.Sprintf("feature-branches = %q\n\n", config.SyncFeatureStrategy))
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.SyncPerennialStrategyHelp)) + "\n")
	result.WriteString(fmt.Sprintf("perennial-branches = %q\n", config.SyncPerennialStrategy))
	return result.String()
}

func Save(config *configdomain.UnvalidatedConfig) error {
	return os.WriteFile(FileName, []byte(RenderTOML(config)), 0o600)
}

func TOMLComment(text string) string {
	if text == "" {
		return ""
	}
	result := []string{}
	for _, line := range strings.Split(text, "\n") {
		if line == "" {
			result = append(result, "#")
		} else {
			result = append(result, "# "+line)
		}
	}
	return strings.Join(result, "\n")
}
