package gohacks

import (
	"fmt"
	"strconv"
	"strings"

	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
)

func ParseBool(text string) (bool, error) {
	switch strings.ToLower(text) {
	case "yes", "on", "enable", "enabled":
		return true, nil
	case "no", "off", "disable", "disabled":
		return false, nil
	}
	return strconv.ParseBool(text)
}

func ParseBoolOpt(text, source string) (Option[bool], error) {
	if text == "" {
		return None[bool](), nil
	}
	parsed, err := ParseBool(text)
	if err != nil {
		return None[bool](), fmt.Errorf(messages.ValueInvalid, source, text)
	}
	return Some(parsed), nil
}
