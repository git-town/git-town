package bitbucketdatacenter

import (
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
)

// type-check to ensure conformance to the Connector interface
var (
	bbdcAPIConnector AuthConnector
	_                forgedomain.Connector = bbdcAPIConnector
)

// AuthConnector provides access to the Bitbucket DataCenter API.
type AuthConnector struct {
	AnonConnector
	log      print.Logger
	token    string
	username string
}
