package configdomain

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/messages"
)

type ConfigFile struct {
	Branches                 Branches                  `toml:"branches"`
	CodeHosting              *CodeHosting              `toml:"code-hosting"`
	SyncStrategy             *SyncStrategy             `toml:"sync-strategy"`
	PushNewbranches          *NewBranchPush            `toml:"push-new-branches"`
	ShipDeleteTrackingBranch *ShipDeleteTrackingBranch `toml:"ship-delete-remote-branch"`
	SyncUpstream             *SyncUpstream             `toml:"sync-upstream"`
}

type Branches struct {
	Main       *domain.LocalBranchName `toml:"main"`
	Perennials domain.LocalBranchNames `toml:"perennials"`
}

type CodeHosting struct {
	Platform       *CodeHostingPlatformName   `toml:"platform"`
	OriginHostname *CodeHostingOriginHostname `toml:"origin-hostname"`
}

type SyncStrategy struct {
	FeatureBranches   *SyncFeatureStrategy   `toml:"feature-branches"`
	PerennialBranches *SyncPerennialStrategy `toml:"perennial-branches"`
}

func (self SyncStrategy) FeatureBranchesOrDefault() SyncFeatureStrategy {
	if self.FeatureBranches == nil {
		return SyncFeatureStrategyMerge
	}
	return *self.FeatureBranches
}

func (self SyncStrategy) PerennialBranchesOrDefault() SyncPerennialStrategy {
	if self.PerennialBranches == nil {
		return SyncPerennialStrategyRebase
	}
	return *self.PerennialBranches
}

func LoadConfigFile() (*ConfigFile, error) {
	file, err := os.Open(".git-branches.toml")
	defer file.Close()
	if err != nil {
		return nil, nil
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf(messages.ConfigFileCannotRead, ".git-branches.yml", err)
	}
	return ParseConfigFile(string(bytes))
}

func EncodeConfigFile(config ConfigFile) string {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(config)
	if err != nil {
		panic(fmt.Sprintf("cannot encode config: %v", err))
	}
	return buf.String()
}

func ParseConfigFile(text string) (*ConfigFile, error) {
	var result ConfigFile
	_, err := toml.Decode(text, &result)
	if err != nil {
		return nil, fmt.Errorf(messages.ConfigFileWrongInput, ".git-branches.yml", err)
	}
	return &result, err
}
