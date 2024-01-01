package configdomain

import (
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v11/src/gohacks"
	"github.com/git-town/git-town/v11/src/messages"
)

// SyncUpstream contains the configuration setting whether to sync with the upstream remote.
type SyncUpstream bool

func (self SyncUpstream) Bool() bool {
	return bool(self)
}

func (self SyncUpstream) String() string {
	return strconv.FormatBool(self.Bool())
}

func NewSyncUpstreamRef(value bool) *SyncUpstream {
	result := SyncUpstream(value)
	return &result
}

func ParseSyncUpstreamRef(value, source string) (*SyncUpstream, error) {
	parsed, err := gohacks.ParseBool(value)
	if err != nil {
		return nil, fmt.Errorf(messages.ValueInvalid, source, value)
	}
	token := SyncUpstream(parsed)
	return &token, nil
}
