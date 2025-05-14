package forgedomain

import (
	"encoding/json"
	"fmt"

	"github.com/git-town/git-town/v20/internal/messages"
)

// SerializableProposal is a wrapper type that makes the Proposal interface serializable to and from JSON.
type SerializableProposal struct {
	Data      Proposal
	ForgeType ForgeType
}

func (self SerializableProposal) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"forge-type": self.ForgeType.String(),
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
	forgeTypeOpt, err := ParseForgeType(forgeTypeName)
	if err != nil {
		return err
	}
	forgeType, hasForgeType := forgeTypeOpt.Get()
	if !hasForgeType {
		return fmt.Errorf(messages.ForgeTypeUnknown, forgeTypeName)
	}
	switch forgeType {
	case ForgeTypeBitbucket:
		var data BitbucketCloudProposalData
		err = json.Unmarshal(mapping["data"], &data)
		self.Data = data
	case ForgeTypeBitbucketDatacenter:
		var data ProposalData
		err = json.Unmarshal(mapping["data"], &data)
		self.Data = Proposal(data)
	case ForgeTypeCodeberg:
		var data ProposalData
		err = json.Unmarshal(mapping["data"], &data)
		self.Data = data
	case ForgeTypeGitHub:
		var data ProposalData
		err = json.Unmarshal(mapping["data"], &data)
		self.Data = data
	case ForgeTypeGitLab:
		var data ProposalData
		err = json.Unmarshal(mapping["data"], &data)
		self.Data = data
	case ForgeTypeGitea:
		var data ProposalData
		err = json.Unmarshal(mapping["data"], &data)
		self.Data = data
	}
	self.ForgeType = forgeType
	return err
}
