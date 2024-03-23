package configdomain

import (
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v13/src/gohacks"
	"github.com/git-town/git-town/v13/src/messages"
)

// Offline is a new-type for the "offline" configuration setting.
type Offline bool

func (offline Offline) Bool() bool {
	return bool(offline)
}

func (offline Offline) String() string {
	return strconv.FormatBool(offline.Bool())
}

func (offline Offline) ToOnline() Online {
	return Online(!offline.Bool())
}

func NewOfflineRef(value, source string) (*Offline, error) {
	boolValue, err := gohacks.ParseBool(value)
	if err != nil {
		return nil, fmt.Errorf(messages.ValueInvalid, source, value)
	}
	token := Offline(boolValue)
	return &token, nil
}

type Online bool

func (online Online) Bool() bool {
	return bool(online)
}
