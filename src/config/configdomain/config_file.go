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

// ConfigFileData is the unvalidated data as read by the TOML parser.
type ConfigFileData struct {
	Branches                 Branches      `toml:"branches"`
	CodeHosting              *CodeHosting  `toml:"code-hosting"`
	SyncStrategy             *SyncStrategy `toml:"sync-strategy"`
	PushNewbranches          *bool         `toml:"push-new-branches"`
	ShipDeleteTrackingBranch *bool         `toml:"ship-delete-remote-branch"`
	SyncUpstream             *bool         `toml:"sync-upstream"`
}

type Branches struct {
	Main       *string  `toml:"main"`
	Perennials []string `toml:"perennials"`
}

type CodeHosting struct {
	Platform       *string `toml:"platform"`
	OriginHostname *string `toml:"origin-hostname"`
}

type SyncStrategy struct {
	FeatureBranches   *string `toml:"feature-branches"`
	PerennialBranches *string `toml:"perennial-branches"`
}

func (self ConfigFileData) Validate() (PartialGitConfig, error) {
	result := PartialGitConfig{}
	if self.Branches.Main != nil {
		result.MainBranch = domain.NewLocalBranchNameRef(*self.Branches.Main)
	}
	if self.Branches.Perennials != nil {
		result.PerennialBranches = domain.NewLocalBranchNamesRef(self.Branches.Perennials...)
	}
	return result, nil
}

// ConfigFile is validated data from the configuration file, ready to be used by the application.
type ConfigFile struct {
	Branches                 Branches                  `toml:"branches"`
	CodeHosting              *CodeHosting              `toml:"code-hosting"`
	SyncStrategy             *SyncStrategy             `toml:"sync-strategy"`
	PushNewbranches          *NewBranchPush            `toml:"push-new-branches"`
	ShipDeleteTrackingBranch *ShipDeleteTrackingBranch `toml:"ship-delete-remote-branch"`
	SyncUpstream             *SyncUpstream             `toml:"sync-upstream"`
}

func EncodeConfigFile(config ConfigFile) string {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(config)
	if err != nil {
		panic(fmt.Sprintf("cannot encode config: %v", err))
	}
	return buf.String()
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
	return ParseConfigFileData(string(bytes))
}

func ParseTOML(text string) (*ConfigFileData, error) {
	var result ConfigFile
	_, err := toml.Decode(text, &result)
	if err != nil {
		return nil, fmt.Errorf(messages.ConfigFileWrongInput, ".git-branches.yml", err)
	}
	return &result, err
}
