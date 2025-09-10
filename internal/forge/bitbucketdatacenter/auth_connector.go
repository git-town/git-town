package bitbucketdatacenter

import (
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// type-check to ensure conformance to the Connector interface
var (
	bbdcAPIConnector AuthConnector
	_                forgedomain.AuthVerifier = bbdcAPIConnector
	_                forgedomain.Connector    = bbdcAPIConnector
)

// AuthConnector provides access to the API of Bitbucket installations.
type AuthConnector struct {
	AnonConnector
	log      print.Logger
	token    string
	username string
}

// TODO: delete this, it doesn't do anything
func (self AuthConnector) VerifyConnection() forgedomain.VerifyConnectionResult {
	return forgedomain.VerifyConnectionResult{
		AuthenticatedUser:   None[string](),
		AuthenticationError: nil,
		AuthorizationError:  nil,
	}
}
