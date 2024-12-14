package configfile

import (
	"fmt"
	"os"
	"strings"

	"github.com/git-town/git-town/v16/internal/cli/dialog"
	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
)

func RenderPerennialBranches(perennials gitdomain.LocalBranchNames) string {
	if len(perennials) == 0 {
		return "[]"
	}
	return fmt.Sprintf(`["%s"]`, perennials.Join(`", "`))
}

func RenderTOML(config *config.UnvalidatedConfig) string {
	result := strings.Builder{}
	result.WriteString("# Git Town configuration file\n")
	result.WriteString("#\n")
	result.WriteString("# Run \"git town config setup\" to add additional entries\n")
	result.WriteString("# to this file after updating Git Town.\n")
	result.WriteString("\n[branches]\n\n")
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.MainBranchHelp)) + "\n")
	result.WriteString(fmt.Sprintf("main = %q\n\n", config.UnvalidatedConfig.MainBranch))
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.PerennialBranchesHelp)) + "\n")
	result.WriteString(fmt.Sprintf("perennials = %s\n", RenderPerennialBranches(config.NormalConfig.PerennialBranches)) + "\n")
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.PerennialRegexHelp)) + "\n")
	result.WriteString(fmt.Sprintf("perennial-regex = %q\n", config.NormalConfig.PerennialRegex))
	result.WriteString("\n[create]\n\n")
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.NewBranchTypeHelp)) + "\n")
	result.WriteString(fmt.Sprintf("new-branch-type = %q\n\n", config.NormalConfig.NewBranchType))
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.PushNewBranchesHelp)) + "\n")
	result.WriteString(fmt.Sprintf("push-new-branches = %t\n", config.NormalConfig.PushNewBranches))
	result.WriteString("\n[hosting]\n\n")
	result.WriteString(fmt.Sprintf("dev-remote = %q\n\n", config.NormalConfig.DevRemote.String()))
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.HostingPlatformHelp)) + "\n")
	if platform, has := config.NormalConfig.HostingPlatform.Get(); has {
		result.WriteString(fmt.Sprintf("platform = %q\n\n", platform))
	} else {
		result.WriteString("# platform = \"\"\n\n")
	}
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.OriginHostnameHelp)) + "\n")
	if config.NormalConfig.HostingOriginHostname.IsNone() {
		result.WriteString("# origin-hostname = \"\"\n")
	} else {
		result.WriteString(fmt.Sprintf("origin-hostname = %q\n", config.NormalConfig.HostingOriginHostname))
	}
	result.WriteString("\n[ship]\n\n")
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.ShipDeleteTrackingBranchHelp)) + "\n")
	result.WriteString(fmt.Sprintf("delete-tracking-branch = %t\n\n", config.NormalConfig.ShipDeleteTrackingBranch))
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.ShipStrategyHelp)) + "\n")
	result.WriteString(fmt.Sprintf("strategy = %q\n", config.NormalConfig.ShipStrategy))
	result.WriteString("\n[sync]\n\n")
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.SyncFeatureStrategyHelp)) + "\n")
	result.WriteString(fmt.Sprintf("feature-strategy = %q\n\n", config.NormalConfig.SyncFeatureStrategy))
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.SyncPerennialStrategyHelp)) + "\n")
	result.WriteString(fmt.Sprintf("perennial-strategy = %q\n\n", config.NormalConfig.SyncPerennialStrategy))
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.SyncPrototypeStrategyHelp)) + "\n")
	result.WriteString(fmt.Sprintf("prototype-strategy = %q\n\n", config.NormalConfig.SyncPrototypeStrategy))
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.PushHookHelp)) + "\n")
	result.WriteString(fmt.Sprintf("push-hook = %t\n\n", config.NormalConfig.PushHook))
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.SyncTagsHelp)) + "\n")
	result.WriteString(fmt.Sprintf("tags = %t\n\n", config.NormalConfig.SyncTags))
	result.WriteString(TOMLComment(strings.TrimSpace(dialog.SyncUpstreamHelp)) + "\n")
	result.WriteString(fmt.Sprintf("upstream = %t\n", config.NormalConfig.SyncUpstream))
	return result.String()
}

func Save(config *config.UnvalidatedConfig) error {
	return os.WriteFile(FileName, []byte(RenderTOML(config)), 0o600)
}

func TOMLComment(text string) string {
	if text == "" {
		return ""
	}
	var result []string
	for _, line := range strings.Split(text, "\n") {
		if line == "" {
			result = append(result, "#")
		} else {
			result = append(result, "# "+line)
		}
	}
	return strings.Join(result, "\n")
}
