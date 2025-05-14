package forgedomain

import (
	"encoding/json"
	"fmt"

	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/messages"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

// Proposal provides information about a change request on a forge.
// Alternative names are "pull request" or "merge request".
type Proposal interface {
	// text of the body of the proposal
	// if Some, the string is guaranteed to be non-empty
	GetBody() Option[string]

	// whether this proposal can be merged via the API
	GetMergeWithAPI() bool

	// the number used to identify the proposal on the forge
	GetNumber() int

	// name of the source branch ("head") of this proposal
	GetSource() gitdomain.LocalBranchName

	// name of the target branch ("base") of this proposal
	GetTarget() gitdomain.LocalBranchName

	// text of the title of the proposal
	GetTitle() string

	// the URL of this proposal
	GetURL() string
}

func CommitBody(proposal Proposal, title string) string {
	result := title
	if body, has := proposal.GetBody().Get(); has {
		result += "\n\n"
		result += body
	}
	return result
}

type ProposalData struct {
	Body         Option[string]
	MergeWithAPI bool
	Number       int
	Source       gitdomain.LocalBranchName
	Target       gitdomain.LocalBranchName
	Title        string
	URL          string
}

func (self ProposalData) GetBody() Option[string] {
	return self.Body
}

func (self ProposalData) GetMergeWithAPI() bool {
	return self.MergeWithAPI
}

func (self ProposalData) GetNumber() int {
	return self.Number
}

func (self ProposalData) GetSource() gitdomain.LocalBranchName {
	return self.Source
}

func (self ProposalData) GetTarget() gitdomain.LocalBranchName {
	return self.Target
}

func (self ProposalData) GetTitle() string {
	return self.Title
}

func (self ProposalData) GetURL() string {
	return self.URL
}

// SerializableProposal is a wrapper type that makes the Proposal interface serializable to and from JSON.
type SerializableProposal struct {
	ForgeType ForgeType
	Data      Proposal
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
		var data ProposalData
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
	return err
}
