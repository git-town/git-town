package bitbucketdatacenter

import (
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// type-check to ensure conformance to the Connector interface
var (
	bbdcAPIConnector APIConnector
	_                forgedomain.APIConnector = bbdcAPIConnector
	_                forgedomain.Connector    = bbdcAPIConnector
)

// APIConnector provides access to the API of Bitbucket installations.
type APIConnector struct {
	WebConnector
	log      print.Logger
	token    string
	username string
}

// TODO: delete this, it doesn't do anything
func (self APIConnector) VerifyConnection() forgedomain.VerifyConnectionResult {
	return forgedomain.VerifyConnectionResult{
		AuthenticatedUser:   None[string](),
		AuthenticationError: nil,
		AuthorizationError:  nil,
	}
}
