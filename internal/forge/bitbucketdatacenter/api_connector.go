package bitbucketdatacenter

import (
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
)

// type-check to ensure conformance to the Connector interface
var (
	bbdcAPIConnector WebConnector
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
