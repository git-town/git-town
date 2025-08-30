package forgedomain

import (
	"encoding/json"
	"fmt"

	"github.com/git-town/git-town/v21/internal/messages"
)

// Proposal is a wrapper type that makes the Proposal interface serializable to and from JSON.
type Proposal struct {
	Data      ProposalInterface
	ForgeType ForgeType
}

func (self Proposal) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"data":       self.Data,
		"forge-type": self.ForgeType.String(),
	})
}

// UnmarshalJSON is used when de-serializing a Proposal from JSON.
func (self *Proposal) UnmarshalJSON(b []byte) error {
	var mapping map[string]json.RawMessage
	if err := json.Unmarshal(b, &mapping); err != nil {
		return err
	}
	var forgeTypeName string
	if err := json.Unmarshal(mapping["forge-type"], &forgeTypeName); err != nil {
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
	case ForgeTypeBitbucketDatacenter, ForgeTypeCodeberg, ForgeTypeGitHub, ForgeTypeGitLab, ForgeTypeGitea:
		var data ProposalData
		err = json.Unmarshal(mapping["data"], &data)
		self.Data = data
	}
	self.ForgeType = forgeType
	return err
}
