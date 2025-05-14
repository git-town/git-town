package forge

import (
	"encoding/json"
	"fmt"

	"github.com/git-town/git-town/v20/internal/forge/bitbucketcloud"
	"github.com/git-town/git-town/v20/internal/forge/bitbucketdatacenter"
	"github.com/git-town/git-town/v20/internal/forge/codeberg"
	"github.com/git-town/git-town/v20/internal/forge/forgedomain"
	"github.com/git-town/git-town/v20/internal/forge/gitea"
	"github.com/git-town/git-town/v20/internal/forge/github"
	"github.com/git-town/git-town/v20/internal/forge/gitlab"
	"github.com/git-town/git-town/v20/internal/messages"
)

// SerializableProposal is a wrapper type that makes the Proposal interface serializable to and from JSON.
type SerializableProposal struct {
	ForgeType forgedomain.ForgeType
	Data      interface{}
}

func (self SerializableProposal) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"forge-type": self.ForgeType,
		"data":       self.Data,
	})
}

// UnmarshalJSON is used when de-serializing JSON into a Location.
func (self *SerializableProposal) UnmarshalJSON(b []byte) error {
	var mapping map[string]json.RawMessage
	err := json.Unmarshal(b, &mapping)
	if err != nil {
		return err
	}
	var forgeTypeName string
	err = json.Unmarshal(mapping["forge-type"], &forgeTypeName)
	if err != nil {
		return err
	}
	forgeTypeOpt, err := forgedomain.ParseForgeType(forgeTypeName)
	if err != nil {
		return err
	}
	forgeType, hasForgeType := forgeTypeOpt.Get()
	if !hasForgeType {
		return fmt.Errorf(messages.ForgeTypeUnknown, forgeTypeName)
	}
	switch forgeType {
	case forgedomain.ForgeTypeBitbucket:
		var data bitbucketcloud.Proposal
		err = json.Unmarshal(mapping["data"], &data)
		self.Data = data
	case forgedomain.ForgeTypeBitbucketDatacenter:
		var data bitbucketdatacenter.Proposal
		err = json.Unmarshal(mapping["data"], &data)
		self.Data = data
	case forgedomain.ForgeTypeCodeberg:
		var data codeberg.Proposal
		err = json.Unmarshal(mapping["data"], &data)
		self.Data = data
	case forgedomain.ForgeTypeGitHub:
		var data github.Proposal
		err = json.Unmarshal(mapping["data"], &data)
		self.Data = data
	case forgedomain.ForgeTypeGitLab:
		var data gitlab.Proposal
		err = json.Unmarshal(mapping["data"], &data)
		self.Data = data
	case forgedomain.ForgeTypeGitea:
		var data gitea.Proposal
		err = json.Unmarshal(mapping["data"], &data)
		self.Data = data
	}
	return err
}
