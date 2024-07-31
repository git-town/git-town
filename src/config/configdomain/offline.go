package configdomain

import (
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v14/src/gohacks"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
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

func NewOfflineOption(valueStr, source string) (Option[Offline], error) {
	if valueStr == "" {
		return None[Offline](), nil
	}
	valueBool, err := gohacks.ParseBool(valueStr)
	if err != nil {
		return None[Offline](), fmt.Errorf(messages.ValueInvalid, source, valueStr)
	}
	return Some(Offline(valueBool)), nil
}

type Online bool

func (online Online) Bool() bool {
	return bool(online)
}
