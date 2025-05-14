package forgedomain

import (
	"encoding/json"
	"fmt"

	"github.com/git-town/git-town/v20/internal/forge/bitbucketcloud"
	"github.com/git-town/git-town/v20/internal/forge/bitbucketdatacenter"
	"github.com/git-town/git-town/v20/internal/forge/codeberg"
	"github.com/git-town/git-town/v20/internal/forge/gitea"
	"github.com/git-town/git-town/v20/internal/forge/github"
	"github.com/git-town/git-town/v20/internal/forge/gitlab"
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/messages"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

// Proposal provides information about a change request on a forge.
// Alternative names are "pull request" or "merge request".
type Proposal interface {
	// text of the body of the proposal
	// if Some, the string is guaranteed to be non-empty
	Body() Option[string]

	// whether this proposal can be merged via the API
	MergeWithAPI() bool

	// the number used to identify the proposal on the forge
	Number() int

	// name of the source branch ("head") of this proposal
	Source() gitdomain.LocalBranchName

	// name of the target branch ("base") of this proposal
	Target() gitdomain.LocalBranchName

	// text of the title of the proposal
	Title() string

	// the URL of this proposal
	URL() string
}

func CommitBody(proposal Proposal, title string) string {
	result := title
	if body, has := proposal.Body().Get(); has {
		result += "\n\n"
		result += body
	}
	return result
}

// SerializableProposal is a wrapper type that makes the Proposal interface serializable to and from JSON.
type SerializableProposal struct {
	ForgeType ForgeType
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
	forgeTypeOpt, err := ParseForgeType(forgeTypeName)
	if err != nil {
		return err
	}
	forgeType, hasForgeType := forgeTypeOpt.Get()
	if !hasForgeType {
		fmt.Errorf(messages.ForgeTypeUnknown, forgeTypeName)
	}
	switch forgeType {
	case ForgeTypeBitbucket:
		var data bitbucketcloud.Proposal
		err = json.Unmarshal(mapping["data"], &data)
		self.Data = data
	case ForgeTypeBitbucketDatacenter:
		var data bitbucketdatacenter.Proposal
		err = json.Unmarshal(mapping["data"], &data)
		self.Data = data
	case ForgeTypeCodeberg:
		var data codeberg.Proposal
		err = json.Unmarshal(mapping["data"], &data)
		self.Data = data
	case ForgeTypeGitHub:
		var data github.Proposal
		err = json.Unmarshal(mapping["data"], &data)
		self.Data = data
	case ForgeTypeGitLab:
		var data gitlab.Proposal
		err = json.Unmarshal(mapping["data"], &data)
		self.Data = data
	case ForgeTypeGitea:
		var data gitea.Proposal
		err = json.Unmarshal(mapping["data"], &data)
		self.Data = data
	}
	return err
}
